package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	"github.com/sheghun/go-grpc/calculator/calculatorpb"
)

const (
	addr = "0.0.0.0:50051"
)

func main() {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	c := calculatorpb.NewCalculateServiceClient(cc)

	//doSumRequest(c)
	doPrimeDecomposition(c)
}

func doSumRequest(c calculatorpb.CalculateServiceClient) {

	req := &calculatorpb.SumRequest{
		Sum: &calculatorpb.Sum{
			A: 3,
			B: 10,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculate RPC")
	}

	log.Printf("Response from Calculate: %v", res.Result)
}

func doPrimeDecomposition(c calculatorpb.CalculateServiceClient) {
	fmt.Printf("starting PrimeDecomposition\n")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC...: %v\n", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Printf("Streaming done\n")
				break
			}

			log.Fatalf("error while creating stream: %v", err)
		}

		log.Printf("Response from PrimeNumberDecomposition: %v", msg.GetResult())
	}
}
