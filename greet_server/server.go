package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/faramarzq/grpc_go_course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.GreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function invoked %v", req)
	firstName := req.GetGreeting().GetFirstName()

	if firstName == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received an invalid argument: %v\n", firstName),
		)
	}

	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	fmt.Println("Received the request")

	for i := 0; i < 10; i++ {
		result := "Hello " + req.GetGreeting().GetFirstName() + " " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function called.")

	result := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// End of client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone is invoked\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "

		time.Sleep(100 * time.Millisecond)

		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Response: result,
		})
		if err != nil {
			log.Fatalf("Error while sending response: %v", err)
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic("err connecting")
	}

	// Create new gRPC server instance
	gRPCServer := grpc.NewServer()

	// Registers the server
	greetpb.RegisterGreetServiceServer(gRPCServer, &server{})

	// serve the server
	err = gRPCServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
