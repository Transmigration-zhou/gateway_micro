package main

import (
	"context"
	"flag"
	"fmt"
	pb "gateway-micro/server_client_demo/grpc/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net"
)

var port = flag.Int("port", 50055, "the port to serve on")

const (
	streamingCount = 10
)

type server struct {
	*pb.UnimplementedEchoServer
}

func (s *server) ServerStreamingEcho(in *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	fmt.Println("--- ServerStreamingEcho ---")
	fmt.Printf("request received: %v\n", in)
	// Read requests and send responses.
	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message %v\n", in.Message)
		err := stream.Send(&pb.EchoResponse{Message: in.Message})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	fmt.Println("--- ClientStreamingEcho ---")
	// Read requests and send responses.
	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("echo last received message\n")
			return stream.SendAndClose(&pb.EchoResponse{Message: message})
		}
		message = in.Message
		fmt.Printf("request received: %v, building echo\n", in)
		if err != nil {
			return err
		}
	}
}

func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	fmt.Println("--- BidirectionalStreamingEcho ---")
	// Read requests and send responses.
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending echo\n", in)
		if err := stream.Send(&pb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("--- UnaryEcho ---\n")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("miss metadata from context")
	}
	fmt.Println("md", md)
	fmt.Printf("request received: %v, sending echo\n", in)
	return &pb.EchoResponse{Message: in.Message}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	s.Serve(lis)
}
