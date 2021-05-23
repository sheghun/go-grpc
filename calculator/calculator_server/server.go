package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/sheghun/go-grpc/calculator/calculatorpb"
)

const (
	port = ":50051"
)

type server struct {
	calculatorpb.UnimplementedCalculateServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v ", err)
	}

	grpcServer := grpc.NewServer()
	calculatorpb.RegisterCalculateServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v ", err)
	}
}

func (s *server) Sum(_ context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function invoked\n")

	a := req.GetSum().GetA()
	b := req.GetSum().GetB()

	sum := a + b

	res := &calculatorpb.SumResponse{
		Result: sum,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculateService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function invoked\n")

	N := req.GetNumber()

	var k uint32 = 2

	for N > 1 {
		time.Sleep(1 * time.Second)
		if N % k == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: k,
			}
			err := stream.Send(res)
			if err != nil {
				fmt.Errorf("error occurred when streaming data %v\n", err)
				return err
			}
			N = N/k
			continue
		}
		k = k + 1
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculateService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage was invoked using RPC...\n")
	var nums float32
	var counter float32
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
					Result: nums/counter,
				})
			}
			log.Fatalf("error while reading client stream %v\n", err)
		}
		nums += float32(req.GetNumber())
		counter++
	}

	return nil
}