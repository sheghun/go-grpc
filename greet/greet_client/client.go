package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/sheghun/go-grpc/greet/greetpb"
)

const (
	addr = "localhost:50051"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect %v: ", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	//doUnary(c)

	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Oladiran",
			LastName:  "Segun",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a Server Streaming RPC...\n")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Oladiran",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling GreatManyTimes RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// We're at the end of the stream of break
				break
			}
			log.Fatalf("error while creating stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")

	reqs := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Segun",
			},
		}, &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Paul",
			},
		}, &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dorcas",
			},
		}, &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deborah",
			},
		}, &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Favour",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	// We itereate over our slice and send each message individually
	for _, req := range reqs {
		fmt.Printf("Send req: %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v\n", err)
	}
	fmt.Printf("LongGreet Response:  %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")
	var wg sync.WaitGroup

	reqs := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Segun",
			},
		}, &greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Paul",
			},
		}, &greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dorcas",
			},
		}, &greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Deborah",
			},
		}, &greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Favour",
			},
		},
	}

	// We create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v\n", err)
		return
	}

	// We send a bunch of messages from the client (go routine)
	wg.Add(1)
	go func() {
		// function to send a bunch of messages
		for _, req := range reqs {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	// We receive a bunch of messages from the client (go routine)
	go func(wg *sync.WaitGroup) {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("Error while receiving: %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		wg.Done()
	}(&wg)

	// block until everything is done
	wg.Wait()
}