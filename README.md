# gRPC with GoLang
This repo hands on lab at 4 model of gRPC
<ol>
  <li><a href="#unary-call">Unary call from client to server</a></li>
  <li><a href="#server-streaming">Client receive streaming from server</a></li>
  <li><a href="#client-streaming">Client stream to server</a></li>
  <li><a href="#bi-directional-streaming">Bi-directional streaming from both client, server</a></li>
  <li><a href="#advanced">Advanced implementation</a>
    <ol>
      <li><a href="#deadline">Deadline</a></li>
      <li><a href="#ssl">SSL</a></li>
    </ol>
  </li>
</ol>

# Service defination
- Protocol buffer define file: [./proto/greet.proto](./proto/greet.proto)
- Makefile: [./makefile](./makefile)
```makefile
create-pb:
	protoc --go_out=. --go-grpc_out=. proto/greet.proto
```
- Run make file to generate grpc for golang

```bash
make create-pb
```
## Define service
```Go
service GreetingService {
  rpc SayHello (HelloRequest) returns (HelloResponse);
  rpc SayHelloServerStream (NameList) returns (stream HelloResponse);
  rpc SayHelloClientStream (stream HelloRequest) returns (MessageList);
  rpc SayHelloBidirectionalStream (stream HelloRequest) returns (stream HelloResponse);
}
```
## Define messages
- 
```go
message HelloRequest {
  string name = 1;
}
```
```go
message HelloResponse {
  string message = 1;
}
```
```go
message NameList {
  repeated string names = 1;
}
```
```go
message MessageList {
  repeated string messages = 1;
}
```
## Genarated pb.gp file
- message struct : [./proto/greet.pb.go](./proto/greet.pb.go)
- gRPC client server : [./proto/greet_grpc.pb.go](./proto/greet_grpc.pb.go)
# Client - Server Implementation

## <a id="unary-call"></a>Unary call from client to server 

#### client
Enter name, wrap it onto ```pb.HelloRequest```, then send it to server by ```pb.GreetingServiceClient```
```go
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
```

#### server
Get sent name from client, log it out, then send back to client by `pb.HelloResponse` 
```go
func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Client named %v is calling", in.GetName())
	return &pb.HelloResponse{
		Message: "Server say hello to " + in.Name,
	}, nil
}
```


## <a id="client-streaming"></a>Client stream to server

#### client
- Create a loop through `pb.NameList` then send continuously to server `pb.HelloRequest` message
```go
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
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not receive response: %v", err)
	}
	log.Printf("Received: %s", res.Messages)
	log.Printf("Client streaming request completed")
}
```

#### server
- Create a loop that listen to client request <br>
- Append each name received to a slice of string until streaming is finished<br>
- Then send `pb.MessageList` back to client
```go
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
```
## <a id="server-streaming"></a> Client receive streaming from server

#### client
- Send to server `pb.NameList` then listen server say hello to **each** name in that list =)))
- That means streaming from server
```go
func callSayHelloServerStream(client pb.GreetingServiceClient, names *pb.NameList) {
	log.Printf("Streaming request started")
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
```

#### server
- After receive a `pb.NameList` as a slice of name from client, slice it then say `pb.HelloResponse` to client each by each name
```go
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
	}
	return nil
}
```
## <a id="bi-directional-streaming"></a> Bi-directional streaming client <-> server

#### client
- Create a channel to streaming listen to server first, then send each name to server, receive hello and push into channel
- The order like:
```
1. create channel waitc to listen hello from server
2. Send 1st name to server, receive hello but push into channel
3. Send 2nd name to server, receive hello but push into channel
4. Send 3rd name to server, receive hello but push into channel
...
10. Send last name to server, receive hello but push into channel
11. Receive EOF, which means nothing streaming
12. Log the output inside the channel to console `<-waitc`, then close channel
```
```go
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
```

#### server
- Whenever receive a nem from client, send hello back to client by `pb.HelloResponse`
```go
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
```
# <a id="advanced"></a>Advanced Section

## <a id="deadline"></a>Deadline
- For example:
  - Client request that server have to say hello within 3 second
<pre><code>
func callSayHelloServerStream(client pb.GreetingServiceClient, names *pb.NameList) {
	log.Printf("Streaming request started")
	<b>ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()</b>
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
</code></pre>
  - But server send each name after 1 second
<pre><code>
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
		<b>time.Sleep(1 * time.Second)</b>
	}
	return nil
}
</code></pre>

- Result from client console
```
1. SayHello
2. SayHelloServerStream
3. SayHelloClientStream
4. SayHelloBidirectionalStream
0. Exit
Choose: 2
2024/03/21 18:01:52 Streaming request started
2024/03/21 18:01:52 Received from server: Hello Pep Guardiola
2024/03/21 18:01:53 Received from server: Hello Kevin De Bruyne
2024/03/21 18:01:54 Received from server: Hello Julian Alvarez
2024/03/21 18:01:55 error while streaming: rpc error: code = DeadlineExceeded desc = context deadline exceeded
exit status 1
```

## <a id="ssl"></a>SSL
