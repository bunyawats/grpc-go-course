package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello Greet Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v\n", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Create client %f\n", c)

	//doUnary(c)
	doServerStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Start to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bunyawat",
			LastName:  "Singchai",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v\n", err)
	}
	fmt.Printf("Response from Greet: %v\n", res.Resulte)

}

func doServerStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Start to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimeRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bunyawat",
			LastName:  "Singchai",
		},
	}
	resStream, err := c.GreetManyTime(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while callong GreetManyTime RPC: %v\n ", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we have reached the end of the straem
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v\n", err)
		}
		fmt.Printf("Resule from GreetManyTime RPC: %v\n", msg.GetResulte())
	}
}
