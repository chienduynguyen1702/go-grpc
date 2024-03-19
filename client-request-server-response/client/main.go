package main

import (
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
	// callSayHello(client)
	callSayHelloServerStream(client, names)
}
