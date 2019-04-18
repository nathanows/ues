package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nathanows/ues/echo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	start := time.Now()

	creds, err := credentials.NewClientTLSFromFile("certs/server-cert.pem", "")
	if err != nil {
		log.Fatalf("cert load error: %s", err)
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	fmt.Println(serverAddr)
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
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

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", os.Getenv("TOKEN"))

	resp, err := client.Echo(ctx, &req)
	if err != nil {
		log.Fatalf("error when calling Echo: %s", err)
	}
	elapsed := time.Since(start)
	ch <- fmt.Sprintf("Response: %+v, Took: %s", resp.Message, elapsed)
}
