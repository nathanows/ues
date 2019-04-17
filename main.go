package main

import (
	"context"
	"log"
	"net"

	echopb "github.com/nathanows/ues/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":6000"

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to start gRPC server: %s", err)
	}

	s := EchoService{}

	grpcServer := grpc.NewServer()
	echopb.RegisterEchoServiceServer(grpcServer, &s)
	reflection.Register(grpcServer)

	log.Printf("starting gRPC server on %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server: %s", err)
	}
}

// EchoService implements gRPC echo.
type EchoService struct {
}

// Echo responds with a message body matching that of the request message body
func (s *EchoService) Echo(ctx context.Context, req *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	return &echopb.EchoResponse{Message: req.Message}, nil
}
