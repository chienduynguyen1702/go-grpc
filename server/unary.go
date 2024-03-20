package main

import (
	"context"
	"log"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
)

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Client named %v is calling", in.GetName())
	return &pb.HelloResponse{
		Message: "Server say hello to " + in.Name,
	}, nil
}
