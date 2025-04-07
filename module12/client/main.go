package main

import (
    "context"
    "fmt"
    "log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "example/hello/module12/routeguide"
)

func main() {
    // Creating client connection with NewClient instead of Dial
    client, err := grpc.NewClient("localhost:50051", 
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("failed to create client: %v", err)
    }
    defer client.Close()
    
    c := pb.NewRouteGuideClient(client)
    point := &pb.Point{Latitude: 409146138, Longitude: -746188906}
    feature, err := c.GetFeature(context.Background(), point)
    if err != nil {
        log.Fatalf("could not get feature: %v", err)
    }
    fmt.Printf("Feature: %s, Location: (%d, %d)\n", feature.Name, feature.Location.Latitude, feature.Location.Longitude)
}