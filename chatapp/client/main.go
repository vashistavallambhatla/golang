package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	pb "example/hello/chatapp/grpc"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatClient(conn)

	reader := bufio.NewReader(os.Stdin)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	log.Print("Enter your username: ")
	sender, _ := reader.ReadString('\n')
	sender = strings.TrimSpace(sender)

	rooms, err := client.GetExistingChatRooms(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Failed to fetch the rooms: %v", err)
	}

	if len(rooms.Rooms) == 0 {
		log.Print("No existing rooms found. You can create one by entering a room name.")
	} else {
		log.Printf("Available rooms: %v.Enter the room you want to join or you can also create your own room by entering a new room name.", rooms.Rooms)
	}

	room, _ := reader.ReadString('\n')
	room = strings.TrimSpace(room)

	roomJoined, err := client.JoinRoom(context.Background(), &pb.JoinRequest{Sender: sender, Room: room})
	if err != nil {
		log.Fatalf("Failed to join room: %v", err)
	}
	log.Printf("You joined room %s with %v other members in the room: %v", room,len(roomJoined.Members),roomJoined.Members)

	go func() {
		<-signalChan
		fmt.Println()
		_, err := client.LeaveChatRoom(context.Background(), &pb.LeaveRequest{
			Sender: sender,
			Room:   room,
			Type:   "sigexit",
		})
		if err != nil {
			log.Printf("Error leaving chat on Ctrl+C: %v", err)
		}
		os.Exit(0)
	}()

	go func() {
		joinReq := &pb.JoinRequest{
			Sender: sender,
			Room:   room,
		}
		stream, err := client.BroadcastRoomUpdate(context.Background(), joinReq)
		if err != nil {
			log.Printf("Error joining update stream: %v", err)
			return
		}
		for {
			update, err := stream.Recv()
			if err != nil {
				log.Printf("Update stream ended: %v", err)
				return
			}
			log.Printf("[UPDATE]: %s", update.Update)
			if update.Type == "joined" && update.Sender == sender {
				fmt.Print("press ENTER to start chatting")
			}
		}
	}()

	stream, err := client.RoomChat(context.Background())
	if err != nil {
		log.Fatalf("Failed to start room chat stream: %v", err)
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			format := "[%s]: %s"
			if msg.Room == "private" {
				format = "Private message from [%s]: %s"
			}
			log.Printf(format, msg.Sender, msg.Content)
		}
	}()

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "/pm") {
			parts := strings.SplitN(text, " ", 3)
			if len(parts) != 3 {
				log.Println("Use: /pm <recipient> <message> to send a private message")
				continue
			}
			recipient := parts[1]
			message := parts[2]

			_, err := client.SendPrivateMessage(context.Background(), &pb.PrivateMessage{
				Sender:    sender,
				Recipient: recipient,
				Content:   message,
			})
			if err != nil {
				log.Printf("Failed to send private message: %v", err)
			}
			continue
		}

		if text == "/exit" {
			_, err := client.LeaveChatRoom(context.Background(), &pb.LeaveRequest{
				Sender: sender,
				Room:   room,
			})
			if err != nil {
				log.Printf("Leave error: %v", err)
			}
			log.Println("You left the room.")
			break
		}

		err := stream.Send(&pb.ChatRoomMessage{
			Sender:  sender,
			Room:    room,
			Content: text,
		})
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			break
		}
	}

	time.Sleep(time.Second)
	stream.CloseSend()
}
