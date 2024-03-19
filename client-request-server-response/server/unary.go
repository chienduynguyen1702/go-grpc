package main

import (
	"context"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
)

func (s *helloServer) SayHello(ctx context.Context, in *pb.NoParam) (*pb.HelloResponse, error) {
	// log.Printf("Received: %v", in.GetNames())
	return &pb.HelloResponse{
		Message: "Hello ",
	}, nil
}
