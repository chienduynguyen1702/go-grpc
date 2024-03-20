package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
)

func callSayHello(client pb.GreetingServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prompt the user to enter their name
	fmt.Println("Enter your name: ")
	var inputName string
	_, err := fmt.Scanf("%s", &inputName)
	if err != nil {
		log.Fatalf("could not read name: %v", err)
	}

	// Call the SayHello method from the client
	res, err := client.SayHello(ctx, &pb.HelloRequest{
		Name: inputName,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("received greeting: %s", res.Message)
}
