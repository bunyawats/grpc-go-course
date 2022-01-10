package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.Greeting.GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Resulte: result,
	}
	return res, nil
}

func (*server) GreetManyTime(req *greetpb.GreetManyTimeRequest, stream greetpb.GreetService_GreetManyTimeServer) error {

	fmt.Printf("GreetManyTime RPC was invocked with %v\n", req)
	firstName := req.Greeting.FirstName
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimeResponse{
			Resulte: result,
		}
		stream.Send(res)
		time.Sleep(time.Second)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {

	fmt.Printf("LongGreet RPC was invocked with stream requested")

	result := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we has finish reading the client
			res := &greetpb.LongGreetResponse{
				Resulte: result,
			}
			return stream.SendAndClose(res)

		}
		if err != nil {
			log.Fatalf("Error while reading client strean: %v\n", err)
		}
		result += "Hello " + req.GetGreeting().GetFirstName() + "!"
	}
}

func (*server) GreetEveryOne(stream greetpb.GreetService_GreetEveryOneServer) error {

	fmt.Printf("GreetEveryOne RPC was invocked with stream requested")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we has finish reading the client
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client strean: %v\n", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "!"
		sendErr := stream.Send(&greetpb.GreetEveryOneResponse{
			Resulte: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v\n", err)
			return sendErr
		}
	}
}

func (*server) GreetWithDeadline(
	ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (
	*greetpb.GreetWithDeadlineResponse, error) {

	fmt.Printf("Greet with Deadline function was invoked with %v\n", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("The client cancel the request!")
			return nil, status.Error(
				codes.Canceled,
				"The client cancel the request!")
		}
		time.Sleep(time.Second)
	}

	firstName := req.Greeting.GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Resulte: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("faild to serve %v\n", err)
	}

}
