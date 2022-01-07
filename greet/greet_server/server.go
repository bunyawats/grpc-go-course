package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
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
