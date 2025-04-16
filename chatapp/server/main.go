package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"errors"

	pb "example/hello/chatapp/grpc"

	"google.golang.org/grpc"
)

type chatServer struct {
    pb.UnimplementedChatServer
    clients     map[string][]chan *pb.ChatRoomMessage // map for rooms and their clients' channels
    activeUsers map[string]chan *pb.PrivateMessage    // map for active user private message channels
    mu          sync.Mutex
}

func NewChatServer() *chatServer {
    return &chatServer{
        clients:    make(map[string][]chan *pb.ChatRoomMessage),
        activeUsers : make(map[string]chan *pb.PrivateMessage),
    }
}

func (s *chatServer) SendPrivateMessage(ctx context.Context, msg *pb.PrivateMessage) (*pb.MessageResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for activeUser, activeUserChan := range s.activeUsers {
		if activeUser == msg.Recipient {
			activeUserChan <- msg

			return &pb.MessageResponse{
				Status: "Message sent",
			}, nil
		}
	}

	return &pb.MessageResponse{
		Status: "Operation failed -- No user found",
	}, errors.New("couldn't send the message -- User not present or might've disconnected")
}

func (s *chatServer) RoomChat(stream pb.Chat_RoomChatServer) error {
	clientChan := make(chan *pb.ChatRoomMessage, 10)
	privateMessageChan := make(chan *pb.PrivateMessage,10)
	var room string

	go func() {
		for msg := range clientChan {
			stream.Send(msg)
		}
	}()

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
			s.mu.Lock()
			s.clients[room] = append(s.clients[room], clientChan)
			s.activeUsers[msg.Sender] = privateMessageChan
			s.broadcastRoomUpdate(fmt.Sprintf("%s has joined the room", msg.Sender), room, true)
			s.mu.Unlock()
		}

		s.mu.Lock()
		for _, ch := range s.clients[room] {
			if ch != clientChan {
				ch <- msg
			}
		}
		s.mu.Unlock()
	}

	
	s.mu.Lock()
	defer s.mu.Unlock()

	
	for i, ch := range s.clients[room] {
		if ch == clientChan {
			s.clients[room] = append(s.clients[room][:i], s.clients[room][i+1:]...)
			s.broadcastRoomUpdate(fmt.Sprintf("%s has left the room", room), room, false)
			break
		}
	}

	s.activeUsers = nil

	close(clientChan)
	close(privateMessageChan)
	return nil
}


func (s *chatServer) broadcastRoomUpdate(content, room string, isJoin bool) {
    update := &pb.ChatRoomMessage{
        Content: content,
        Room:    room,
        IsJoin:  isJoin,
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    for _, ch := range s.clients[room] {
        ch <- update
    }
}

func (s *chatServer) LeaveRoom(ctx context.Context, req *pb.LeaveRequest) (*pb.MessageResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	clientChans := s.clients[req.Room]
	for _, ch := range clientChans {
		// no direct way to identify the sender's channel, so you'd need a better mapping or metadata
		// this part depends on how you associate senders to channels in a real app
		_ = ch // placeholder
	}

	// Just remove sender from activeUsers map
	delete(s.activeUsers, req.Sender)

	s.broadcastRoomUpdate(fmt.Sprintf("%s has left the room", req.Sender), req.Room, false)

	return &pb.MessageResponse{Status: "Left the room"}, nil
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
