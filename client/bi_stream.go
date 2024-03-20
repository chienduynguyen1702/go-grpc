package main

import (
	"context"
	"io"
	"log"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
)

func callSayHelloBidirectionStream(client pb.GreetingServiceClient, names *pb.NameList) {
	log.Printf("Bidirectional streaming request started")

	stream, err := client.SayHelloBidirectionalStream(context.Background())
	if err != nil {
		log.Fatalf("could not create streaming client: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					log.Fatalf("Failed to receive a note : %v", err)
				}
				break // Exit the loop if EOF is received
			}
			log.Printf("Received from server : %s", res.Message)
		}
		close(waitc)
	}()
	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
		log.Printf("Calling to server by : %s", name)
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("Failed to close send: %v", err)
	}
	<-waitc
	log.Printf("Bidirectional streaming request completed")
}
