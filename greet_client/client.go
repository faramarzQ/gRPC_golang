package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/faramarzq/grpc_go_course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	// withInsecure has no ssl, shouldn't be on production
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	// doUnaryRPC(c)

	doServerStreaming(c)
}

func doUnaryRPC(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "faramarz",
			LastName:  "qoshchi",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greeing: %v", err)
	}

	log.Printf("Response: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetManyTimesRequest{
		Greeing: &greetpb.Greeting{
			FirstName: "Faramarz",
			LastName:  "Qoshchi",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greeing: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		fmt.Println(msg.GetResult())
	}
}
