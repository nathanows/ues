package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/nathanows/ues/echo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const grpcPort = "localhost:6000"

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to start gRPC server - err=%s", err)
	}

	s := echo.EchoService{}

	creds, err := credentials.NewServerTLSFromFile("certs/server-cert.pem", "certs/server-key.pem")
	if err != nil {
		log.Fatalf("unable to load TLS key - err=%s", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds), withLogUnaryInterceptor()}

	grpcServer := grpc.NewServer(opts...)
	echo.RegisterEchoServiceServer(grpcServer, &s)
	reflection.Register(grpcServer)

	log.Printf("starting gRPC server - port=%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server - err=%s", err)
	}
}

func withLogUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(loggingInterceptor)
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	log.Printf("request processed - method=%s duration=%s error=%v\n", info.FullMethod, time.Since(start), err)
	return h, err
}
