package main

import (
	context "context"
	"fmt"
	"grpc-basics/greetpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type greetServiceServer struct {
}

func (s *greetServiceServer) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("\nRequest long greet")
	for  {
		msg, err := stream.Recv()
		if err == io.EOF {
			_ = stream.SendAndClose(&greetpb.LongGreetResponse{Result: "end"})
			break
		}
		if err != nil {
			log.Fatalf("\nError while reading client stream: %v\n", err)
		}

		fmt.Printf("\nHello long %v\n", msg.Greeting.FirstName)
	}
	return nil
}

func (s *greetServiceServer) GreetManyTimes(request *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fisrtName := request.Greeting.FirstName
	fmt.Printf("\nRequest first name %v", fisrtName)
	for i := 1; i <= 10; i++ {
		result := fmt.Sprintf("Hello time %v for %v", i, fisrtName)
		res := &greetpb.GreetManyTimesResponse{Result: result}
		err := stream.Send(res)
		if err != nil {
			log.Printf("\nError sending response to client: %v", err)
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (s *greetServiceServer) Greet(ctx context.Context, request *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fisrtName := request.Greeting.FirstName
	fmt.Printf("\nRequest first name %v", fisrtName)
	result := fmt.Sprintf("Hello %v", fisrtName)
	res := &greetpb.GreetResponse{Result: result}
	return res, nil
}

func main() {
	fmt.Println("Configuring grpc server...")
	address := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Printf("\nGreeting server listening on %v\n", address)
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &greetServiceServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
