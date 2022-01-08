package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"

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

	//doUnary(c)

	//doServerStreaming(c)

	doClientStreaming(c)
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

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Start to do a Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12345,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while callong PrimeNumberDecomposition RPC: %v\n ", err)
	}
	for {
		result, err := resStream.Recv()
		if err == io.EOF {
			// we have reached the end of the straem
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v\n", err)
		}
		fmt.Printf("Resule from PrimeNumberDecomposition RPC: %v\n", result.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Start to do a Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage RPC: %v\n ", err)
	}

	numbers := []int64{2, 4, 5, 6, 7, 8, 9, 9, 100}

	// we iterate over message sloce and send it individually
	for _, number := range numbers {
		fmt.Printf("Sening req: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
		time.Sleep(100 * time.Microsecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Erroe while rechive response from LongGreet RPC: %v\n", err)
	}
	fmt.Printf("The average is %v\n", res.GetAverage())
}
