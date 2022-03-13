package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/faramarzq/grpc_go_course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// withInsecure has no ssl, shouldn't be on production
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	doUnaryRPC(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBidirectionalStreaming(c)
}

func doUnaryRPC(c greetpb.GreetServiceClient) {

	var firstName string
	if len(os.Args) >= 2 {
		firstName = os.Args[1]
	} else {
		firstName = ""
	}
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: firstName,
			LastName:  "qoshchi",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			fmt.Printf(resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("You've sent no arguments, pass the name as the second arg in cmd.")
				return
			}
		} else {
			log.Fatalf("error while calling greeting: %v", err)
			return
		}
	}

	log.Printf("Response: %v\n", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Faramarz",
			LastName:  "Qoshchi",
		},
	}

	// Send greeting
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greeting: %v", err)
	}

	// Receive greeting stream
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	reqs := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "One",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Two",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Three",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Four",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Five",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling long greet: %v", err)
	}

	for _, req := range reqs {
		fmt.Printf("Sending Req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet response: %v\n", res)
}

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v\n", err)
	}

	reqs := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "One",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Two",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Three",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Four",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Five",
			},
		},
	}

	wait := make(chan string)

	go func() {
		for _, req := range reqs {
			fmt.Printf("Sending request: %v\n", req.GetGreeting().FirstName)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResponse())
		}
		close(wait)
	}()

	<-wait
}
