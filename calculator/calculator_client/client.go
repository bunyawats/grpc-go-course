package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	//doClientStreaming(c)

	//doBiDiStreaming(c)

	doErrorUnary(c)
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

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Start to do a BiDi Streaming RPC...")

	//we create stream by invoke the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while calling FindMaximum RPC: %v\n ", err)
		return
	}

	waitc := make(chan struct{})

	numbers := []int32{2, 3, 3, 5, 30, 10, 45, 22}

	go func() {
		for _, number := range numbers {
			fmt.Printf("Send number: %v\n", number)
			_ = stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(time.Second)
		}
	}()

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
			fmt.Printf("Maximum number is %v\n", res.GetMaximum())

		}
		close(waitc)
	}()

	<-waitc

}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {

	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -10)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {

	fmt.Println("Call doErrorUnary RPC...")

	res, err := c.SquareRoot(context.Background(),
		&calculatorpb.SquareRootRequest{Number: number})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code().String())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Printf("We probably sent nagative number\n")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v\n", err)
		}
	} else {
		fmt.Printf("Result of square root of %v is %v", number, res.GetNumberRoot())
	}

}
