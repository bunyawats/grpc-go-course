package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	fmt.Printf("Received Sum RPC: %v\n", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecoundNumber

	sumResult := firstNumber + secondNumber

	res := &calculatorpb.SumResponse{
		SumResult: sumResult,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {

	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)
	number := req.Number
	divisor := int64(2)

	for number > 1 {

		if number%divisor == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			}
			stream.Send(res)
			time.Sleep(time.Second)
			number /= divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to: %v\n", divisor)
		}

	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {

	fmt.Printf("ComputeAverage RPC was invocked with stream requested\n")

	var sum, count int64

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we has finish reading the client
			avarage := float64(sum) / float64(count)
			res := &calculatorpb.ComputeAverageResponse{
				Average: avarage,
			}
			return stream.SendAndClose(res)

		}
		if err != nil {
			log.Fatalf("Error while reading client strean: %v\n", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("faild to serve %v\n", err)
	}

}
