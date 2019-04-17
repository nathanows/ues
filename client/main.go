package main

import (
	"context"
	"log"
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

	req := echo.EchoRequest{Message: "something"}
	log.Printf("Request: %#v", req)
	response, err := c.Echo(context.Background(), &req)
	if err != nil {
		log.Fatalf("error when calling Echo: %s", err)
	}
	elapsed := time.Since(start)
	log.Printf("Response: %#v, Took: %s", response, elapsed)
}
