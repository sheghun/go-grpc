package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

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
	//doPrimeDecomposition(c)
	//doComputeAverage(c)
	doFindMaximum(c)
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

func doComputeAverage(c calculatorpb.CalculateServiceClient) {
	fmt.Println("starting to do a client streaming")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage %v\n", err)
	}

	for i := 1; i <= 4; i++ {
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: uint32(i),
		})
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from server: %v\n", err)
	}

	fmt.Printf("ComputeAverage Response: %v\n", res)
}

func doFindMaximum(c calculatorpb.CalculateServiceClient) {
	fmt.Println("starting to do a client streaming")
	var wg sync.WaitGroup

	numbers := []int32{
		1,
		5,
		3,
		6,
		2,
		20,
	}

	results := []int32{}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error whille calling FindMaximum: %v", err)
	}

	wg.Add(1)
	go func() {
		for _, num := range numbers {
			fmt.Printf("Sending message: %v\n", &calculatorpb.FindMaximumRequest{Number: num})
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: num,
			})
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("error while receiving response: %v", err)
			}
			fmt.Printf("Received: %v\n", res)
			results = append(results, res.GetResult())
		}
		wg.Done()
	}()

	wg.Wait()

	fmt.Printf("Results: %v", results)

}
