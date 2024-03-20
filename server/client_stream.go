package main

import (
	"io"
	"log"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
)

func (s *helloServer) SayHelloClientStream(stream pb.GreetingService_SayHelloClientStreamServer) error {
	var nameList []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Sent Response Hello Message to Client")
			return stream.SendAndClose(&pb.MessageList{Messages: nameList})
		}
		if err != nil {
			return err
		}
		log.Printf("Streaming from client: %s is calling", req.Name)
		nameList = append(nameList, "Hello ", req.Name, "\n")
	}
}
