package main

import (
	"log"
	"time"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
)

func (s *helloServer) SayHelloServerStream(req *pb.NameList, stream pb.GreetingService_SayHelloServerStreamServer) error {
	log.Fatalf("Received: %v", req.Names)
	for _, name := range req.Names {
		res := &pb.HelloResponse{
			Message: "Hello " + name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
