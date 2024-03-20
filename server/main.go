package main

import (
	"log"
	"net"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
	"google.golang.org/grpc"
)

const (
	// serverAddress is the address of the server.
	host = ""
	port = ":8080"
)

type helloServer struct {
	pb.GreetingServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreetingServiceServer(grpcServer, &helloServer{})
	log.Printf("Server listening at %s", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
