package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "example/hello/chatapp/grpc"

	"google.golang.org/grpc"
)

type chatServer struct {
	pb.UnimplementedChatServer
	rooms     map[string]map[string]chan *pb.ChatRoomMessage
	activeUsers map[string]chan *pb.PrivateMessage
	mu          sync.Mutex
	updates     map[string]map[string]chan *pb.Update
}

func NewChatServer() *chatServer {
	return &chatServer{
		rooms:     make(map[string]map[string]chan *pb.ChatRoomMessage),
		activeUsers: make(map[string]chan *pb.PrivateMessage),
		updates:     make(map[string]map[string]chan *pb.Update),
	}
}

func (s *chatServer) GetAvailableRooms(ctx context.Context, _ *pb.Empty) (*pb.AvailableRooms, error) {
	var rooms []string
	for room := range s.rooms {
		rooms = append(rooms, room)
	}
	return &pb.AvailableRooms{Rooms: rooms}, nil
}

func (s *chatServer) JoinRoom(ctx context.Context, joinReq *pb.JoinRequest) (*pb.JoinRoomResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room := joinReq.Room
	sender := joinReq.Sender

	fmt.Println(sender)
	fmt.Println(s.rooms)

	var users []string

	if s.rooms[room] == nil {
		s.rooms[room] = make(map[string]chan *pb.ChatRoomMessage)
	} else {
		for user := range s.rooms[room] {
			users = append(users, user)
		}
	}
	
	if _, exists := s.rooms[room][sender]; exists {
		return &pb.JoinRoomResponse{
			Status:  "Failed",
			Members: users,
		}, errors.New("username already taken in this room")
	}

	if _, exists := s.rooms[room][sender]; exists {
		return &pb.JoinRoomResponse{
			Status: "Failed",
			Members : users,
		}, errors.New("username already taken in this room")
	}

	s.rooms[room][sender] = make(chan *pb.ChatRoomMessage, 10)
	s.activeUsers[sender] = make(chan *pb.PrivateMessage, 10)

	if s.updates[room] == nil {
		s.updates[room] = make(map[string]chan *pb.Update)
	}
	s.updates[room][sender] = make(chan *pb.Update, 10)

	s.notifyRoomUpdate(room, sender, "joined")

	return &pb.JoinRoomResponse{
		Status: "Success",
		Members: users,
	}, nil
}

func (s *chatServer) notifyRoomUpdate(room, user, action string) {
	update := &pb.Update{
		Sender: user,
		Room:   room,
		Update: user + " has " + action + " the room",
		Type:   action,
	}

	if roomUsers, ok := s.updates[room]; ok {
		for _, updateChan := range roomUsers {
			select {
			case updateChan <- update:
			default:
				log.Printf("Failed to send update to a user: buffer full")
			}
		}
	}
}

func (s *chatServer) BroadcastRoomUpdate(req *pb.JoinRequest, stream pb.Chat_BroadcastRoomUpdateServer) error {
	room := req.Room
	sender := req.Sender

	s.mu.Lock()
	userChan, ok := s.updates[room][sender]
	s.mu.Unlock()

	if !ok {
		return errors.New("user not registered for updates")
	}

	for update := range userChan {
		if update.Type != "sigexit" {
			if err := stream.Send(update); err != nil {
				log.Printf("error sending update to %s: %v", sender, err)
				break
			}
		}
	}

	return nil
}

func (s *chatServer) LeaveChatRoom(ctx context.Context, leaveReq *pb.LeaveRequest) (*pb.MessageResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room := leaveReq.Room
	sender := leaveReq.Sender
	leaveType := leaveReq.Type

	if _, exists := s.rooms[room]; !exists {
		return &pb.MessageResponse{Status: "Failed"}, errors.New("room does not exist")
	}

	if _, exists := s.rooms[room][sender]; !exists {
		return &pb.MessageResponse{Status: "Failed"}, errors.New("user not in the room")
	}

	s.notifyRoomUpdate(room, sender, leaveType)

	if ch := s.rooms[room][sender]; ch != nil {
		close(ch)
	}

	delete(s.rooms[room], sender)

	userInOtherRooms := false
	for r, users := range s.rooms {
		if r != room && users[sender] != nil {
			userInOtherRooms = true
			break
		}
	}

	if !userInOtherRooms {
		if pmChan, exists := s.activeUsers[sender]; exists {
			close(pmChan)
			delete(s.activeUsers, sender)
		}
	}

	if s.updates[room] != nil {
		if updateChan := s.updates[room][sender]; updateChan != nil {
			close(updateChan)
		}
		delete(s.updates[room], sender)
	}

	return &pb.MessageResponse{
		Status: "User left the room successfully",
	}, nil
}

func (s *chatServer) SendPrivateMessage(ctx context.Context, msg *pb.PrivateMessage) (*pb.MessageResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	recipientChan, exists := s.activeUsers[msg.Recipient]
	if !exists {
		return &pb.MessageResponse{
			Status: "Operation failed -- No user found",
		}, errors.New("couldn't send the message -- User not present or might've disconnected")
	}

	select {
	case recipientChan <- msg:
		return &pb.MessageResponse{
			Status: "Message sent",
		}, nil
	default:
		return &pb.MessageResponse{
			Status: "Operation failed -- Message buffer full",
		}, errors.New("couldn't send the message -- Recipient message buffer is full")
	}
}

func (s *chatServer) RoomChat(stream pb.Chat_RoomChatServer) error {
	var sender, room string
	clientChan := make(chan *pb.ChatRoomMessage, 10) // So when a user calls RoomChat() from their end, the stream instance is unique to them and the associated chan is also unique?
	privateMessageChan := make(chan *pb.PrivateMessage, 10)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("recv error", err)
				break
			}

			if room == "" {
				room = msg.Room
				sender = msg.Sender

				s.mu.Lock()
				if s.rooms[room] == nil {
					s.rooms[room] = make(map[string]chan *pb.ChatRoomMessage)
				}
				s.rooms[room][sender] = clientChan
				s.activeUsers[sender] = privateMessageChan
				s.mu.Unlock()

				continue
			}

			s.mu.Lock()
			for user, ch := range s.rooms[room] {
				if user != msg.Sender && ch != nil {

					select {
					case ch <- msg:
					default:
						log.Printf("Failed to send message to %s: buffer full", user)
					}
				}
			}
			s.mu.Unlock()
		}

		s.mu.Lock()
		if s.rooms[room] != nil {
			delete(s.rooms[room], sender)
			if len(s.rooms[room]) == 0 {
				delete(s.rooms, room)
			}
		}
		s.notifyRoomUpdate(room, sender, "left")
		s.mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		for {
			select {
			case roomMsg, ok := <-clientChan:
				if !ok {
					return
				}
				if err := stream.Send(roomMsg); err != nil {
					log.Printf("Error sending room message: %v", err)
					return
				}
			case privateMsg, ok := <-privateMessageChan:
				if !ok {
					return
				}

				roomMsg := &pb.ChatRoomMessage{
					Sender:  privateMsg.Sender,
					Room:    "private",
					Content: 
					privateMsg.Content,
				}
				if err := stream.Send(roomMsg); err != nil {
					log.Printf("Error sending private message: %v", err)
					return
				}
			}
		}
	}()

	wg.Wait()

	s.mu.Lock()
	delete(s.activeUsers, sender)
	if s.updates[room] != nil {
		if updateCh := s.updates[room][sender]; updateCh != nil {
			close(updateCh)
		}
		delete(s.updates[room], sender)
	}
	s.mu.Unlock()

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	chatSrv := NewChatServer()

	pb.RegisterChatServer(grpcServer, chatSrv)
	log.Println("Server is listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
