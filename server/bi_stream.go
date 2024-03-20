package main

import (
	"io"
	"log"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
)

func (s *helloServer) SayHelloBidirectionalStream(stream pb.GreetingService_SayHelloBidirectionalStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Calling from client : %s", req.Name)
		if err = stream.Send(&pb.HelloResponse{
			Message: "Server say hello to " + req.Name,
		}); err != nil {
			return err
		}
	}
}
