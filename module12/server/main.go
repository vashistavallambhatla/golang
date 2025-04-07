package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "log"
    "net"

	pb "example/hello/module12/routeguide"
)

type server struct {
	pb.UnimplementedRouteGuideServer
}

func (s *server)  GetFeature(ctx context.Context,point *pb.Point) (*pb.Feature,error) {
	if point.Latitude == 0 && point.Longitude == 0 {
		return nil,fmt.Errorf("Invalid coordinates")
	}
	return &pb.Feature{Name: "Feature name", Location: point},nil
}

func main() {
	lis, err := net.Listen("tcp",":50051")
	if err!=nil {
		log.Fatalf("Failed to listen: %v",err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer,&server{})
	reflection.Register(grpcServer)

	log.Println("Server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v",err)
	}
}
