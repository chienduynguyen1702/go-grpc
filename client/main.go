package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port   = "8080"
	server = "localhost"
	target = server + ":" + port
)

func main() {

	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}
	var creds credentials.TransportCredentials
	// Check if SSL_MODE is true, then enable SSL
	sslMode := os.Getenv("SSL_MODE")
	if sslMode == "true" {
		certFileServer := "ssl/server.crt"
		sslCreds, sslErr := credentials.NewClientTLSFromFile(certFileServer, server)
		if sslErr != nil {
			log.Fatalf("Failed to generate credentials: %v", sslErr)
		}
		creds = sslCreds
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreetingServiceClient(conn)

	names := &pb.NameList{
		Names: []string{"Pep Guardiola", "Kevin De Bruyne", "Julian Alvarez", "Errling Haland", "Bernardo Silva", "Phil Folden", "John Stones", "Kyle Walker", "Rodri Rodrigo", "Ederson Moraes", "Jeremy Doku"},
	}
	menu := -1
menuLoop:
	for menu != 0 {
		fmt.Printf("1. SayHello\n2. SayHelloServerStream\n3. SayHelloClientStream\n4. SayHelloBidirectionalStream\n0. Exit\n")
		fmt.Printf("Choose: ")
		fmt.Scanf("%d\n", &menu)
		switch menu {
		case 1:
			callSayHello(client)
		case 2:
			callSayHelloServerStream(client, names)
		case 3:
			callSayHelloClientStream(client, names)
		case 4:
			callSayHelloBidirectionStream(client, names)
		case 0:
			break menuLoop
		default:
			fmt.Println("Invalid option")
		}
	}
}

func callSayHello(client pb.GreetingServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prompt the user to enter their name
	fmt.Printf("Enter your name: ")
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

func callSayHelloServerStream(client pb.GreetingServiceClient, names *pb.NameList) {
	log.Printf("Streaming request started")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// ctx := context.Background()
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
		// time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not receive response: %v", err)
	}
	log.Printf("Received: %s", res.Messages)
	log.Printf("Client streaming request completed")
}

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
