package main

import (
	"context"
	"fmt"
	"log"

	pb "example/hello/chatapp/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new chat client
	client := pb.NewChatClient(conn)

	// Get username and room from the user
	fmt.Println("Enter your username:")
	var username string
	fmt.Scanln(&username)

	fmt.Println("Enter the room you want to join:")
	var room string
	fmt.Scanln(&room)

	// Start a new chat stream
	stream, err := client.RoomChat(context.Background())
	if err != nil {
		log.Fatalf("could not create stream: %v", err)
	}

	// Send join message to the room
	joinMessage := &pb.ChatRoomMessage{
		Sender: username,
		Room:   room,
		Content: fmt.Sprintf("%s has joined the room", username),
		IsJoin: true,
	}
	err = stream.Send(joinMessage)
	if err != nil {
		log.Fatalf("failed to send join message: %v", err)
	}
	fmt.Println("Joined the room:", room)

	// Run a goroutine to receive messages from the server
	go receiveMessages(stream)

	// Send and receive messages interactively
	for {
		fmt.Println("Enter message to send (or type 'exit' to leave):")
		var message string
		fmt.Scanln(&message)

		if message == "exit" {
			// Send leave message to the room
			leaveMessage := &pb.ChatRoomMessage{
				Sender: username,
				Room:   room,
				Content: fmt.Sprintf("%s has left the room", username),
				IsJoin: false,
			}
			err = stream.Send(leaveMessage)
			if err != nil {
				log.Fatalf("failed to send leave message: %v", err)
			}

			// Send leave room request to the server
			_, err := client.LeaveRoom(context.Background(), &pb.LeaveRoomRequest{
				Username: username,
				Room:     room,
			})
			if err != nil {
				log.Fatalf("could not leave room: %v", err)
			}
			fmt.Println("You have left the room.")
			break
		}

		// Send the message to the room
		err = stream.Send(&pb.ChatRoomMessage{
			Sender:  username,
			Room:    room,
			Content: message,
			IsJoin:  false,
		})
		if err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
	}
}

// Function to receive and display incoming messages from the server
func receiveMessages(stream pb.Chat_RoomChatClient) {
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving message: %v", err)
		}
		fmt.Printf("\n[%s] %s: %s\n", msg.Room, msg.Sender, msg.Content)
	}
}
