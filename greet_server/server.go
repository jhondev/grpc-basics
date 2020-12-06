package main

import (
	context "context"
	"fmt"
	"grpc-basics/greetpb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type greetServiceServer struct {
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
