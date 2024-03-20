package main

import (
	"context"
	"log"
	"time"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
)

func callSayHelloClientStream(client pb.GreetingServiceClient, names *pb.NameList) {
	log.Printf("Client streaming request started")

	stream, err := client.SayHelloClientStream(context.Background())
	if err != nil {
		log.Fatalf("could not create streaming client: %v", err)
	}
	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("could not send name: %v", err)
		}
		log.Printf("Calling to server by : %s", name)
		time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not receive response: %v", err)
	}
	log.Printf("Received: %s", res.Messages)
	log.Printf("Client streaming request completed")
}
