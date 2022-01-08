package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"

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
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
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

func doClientStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Start to do a Client Streaming RPC...")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet RPC: %v\n ", err)
	}

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
			},
		},
	}

	// we iterate over message sloce and send it individually
	for _, req := range requests {
		fmt.Printf("Sening req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Microsecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while rechive response from LongGreet RPC: %v\n", err)
	}
	fmt.Printf("LongGreet response: %v\n", res.GetResulte())
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Start to do a BiDi Streaming RPC...")

	//we create stream by invoke the client
	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryOne RPC: %v\n ", err)
		return
	}

	waitc := make(chan struct{})

	requests := []*greetpb.GreetEveryOneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
			},
		},
	}

	//we send a bounch of message to the client (go routine)
	go func() {

		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	//we receive a bounch of message from the client (ge routine)
	go func() {

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				//we has reached the last server message
				break
			}
			if err != nil {
				log.Fatalf("Error while reading the message from stream: %v\n", err)
				break
			}
			fmt.Printf("Received %v\n", res.GetResulte())
		}

		close(waitc)
	}()

	//block until everting is done
	<-waitc

}
