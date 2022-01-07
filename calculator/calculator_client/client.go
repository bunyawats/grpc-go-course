package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v\n", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// fmt.Printf("Create client %f\n", c)

	doUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Start to do a Unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:   20,
		SecoundNumber: 30,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v\n", err)
	}
	fmt.Printf("Response from Sum: %v\n", res.SumResult)

}
