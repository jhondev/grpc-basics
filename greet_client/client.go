package main

import (
	"context"
	"fmt"
	"grpc-basics/greetpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Configuring grpc client...")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)
	// doUnary(client)
	doServerStreaming(client)
}

func doUnary(client greetpb.GreetServiceClient) {
	fmt.Println("Requesting Unary...")
	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "jhondev", LastName: "m"}}
	resp, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error in the api: %v", err)
	}
	fmt.Printf("Greeting response: %v", resp.Result)
}

func doServerStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("Requesting stream server...")
	req := &greetpb.GreetManyTimesRequest{Greeting: &greetpb.Greeting{FirstName: "jhondev", LastName: "m"}}
	stream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error in the api: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break // end of stream
		}
		if err != nil {
			log.Fatalf("\nerror receiving server message: %v\n", msg)
		}

		fmt.Printf("\nGreeting many times response: %v\n", msg.Result)
	}
}
