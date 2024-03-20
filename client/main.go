package main

import (
	"fmt"
	"log"

	pb "github.com/chienduynguyen1702/go-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = "8080"
)

func main() {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
