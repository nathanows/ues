package main

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/nathanows/ues/echo"

	"google.golang.org/grpc"
)

func testServer() {
	lis, err := net.Listen("tcp", ":60001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	echo.RegisterEchoServiceServer(s, &echo.Service{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func TestMain(m *testing.M) {
	go testServer()
	os.Exit(m.Run())
}

func TestEcho(t *testing.T) {
	const address = "localhost:60001"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := echo.NewEchoServiceClient(conn)

	t.Run("Echo", func(t *testing.T) {
		request := "Live long and prosper"
		r, err := c.Echo(context.Background(), &echo.EchoRequest{Message: request})
		if err != nil {
			t.Fatalf("could not echo: %v", err)
		}
		if r.Message != request {
			t.Errorf("Expected: %s, got: %s", request, r.Message)
		}

	})
}
