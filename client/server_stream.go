package main

import (
	"context"
	"io"
	"log"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
)

func callSayHelloServerStream(client pb.GreetingServiceClient, names *pb.NameList) {
	log.Printf("Streaming request started")
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	ctx := context.Background()
	stream, err := client.SayHelloServerStream(ctx, names)
	if err != nil {
		log.Fatalf("could not send names: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while streaming: %v", err)
		}
		log.Printf("Received from server: %s", res.Message)
	}
	log.Printf("Streaming request completed")
}
