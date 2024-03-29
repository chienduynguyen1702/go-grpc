package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/chienduynguyen1702/go-grpc/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	// Listen on port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	// define a gRPC server attibutes
	var grpcOptions grpc.ServerOption
	var grpcServer *grpc.Server

	// Load .env file
	envFile, err := godotenv.Read(".env")
	if err != nil {
		panic(err.Error())
	}
	// Check if SSL_MODE is true, then enable SSL
	sslMode := envFile["SSL_MODE"]
	if sslMode == "true" {
		certFileServer := "ssl/server.crt"
		keyFileServer := "ssl/server.key"
		creds, sslErr := credentials.NewServerTLSFromFile(certFileServer, keyFileServer)
		if sslErr != nil {
			log.Fatalf("Failed to generate credentials: %v", sslErr)
		}
		grpcOptions = grpc.Creds(creds)
	}

	// Create a gRPC server object
	grpcServer = grpc.NewServer(grpcOptions)

	// Register the service with the server
	pb.RegisterGreetingServiceServer(grpcServer, &helloServer{})
	log.Printf("Server listening at %s", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Client named %v is calling", in.GetName())
	return &pb.HelloResponse{
		Message: "Server say hello to " + in.Name,
	}, nil
}

func (s *helloServer) SayHelloServerStream(req *pb.NameList, stream pb.GreetingService_SayHelloServerStreamServer) error {
	log.Printf("Received list of name : %v", req.Names)
	for _, name := range req.Names {
		log.Printf("Saying hello to %s", name)
		res := &pb.HelloResponse{
			Message: "Hello " + name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

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
		log.Printf("Sent hello %s to client", req.Name)
	}
}
