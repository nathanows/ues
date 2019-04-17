package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nathanows/ues/echo"

	"google.golang.org/grpc"
)

const grpcAddr = "localhost:6000"

func main() {
	start := time.Now()

	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}
	defer conn.Close()

	c := echo.NewEchoServiceClient(conn)

	ch := make(chan string)
	for _, msg := range os.Args[1:] {
		go makeRequest(msg, c, ch)
	}

	for range os.Args[1:] {
		log.Printf(<-ch)
	}

	elapsed := time.Since(start)
	log.Printf("Completed %d requests in %s", len(os.Args[1:]), elapsed)
}

func makeRequest(msg string, client echo.EchoServiceClient, ch chan<- string) {
	start := time.Now()
	req := echo.EchoRequest{Message: msg}
	log.Printf("Request: %s", req.Message)
	resp, err := client.Echo(context.Background(), &req)
	if err != nil {
		log.Fatalf("error when calling Echo: %s", err)
	}
	elapsed := time.Since(start)
	ch <- fmt.Sprintf("Response: %+v, Took: %s", resp.Message, elapsed)
}
