package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/nathanows/ues/echo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":6000"

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to start gRPC server: %s", err)
	}

	s := echo.EchoService{}

	grpcServer := grpc.NewServer(withLogUnaryInterceptor())
	echo.RegisterEchoServiceServer(grpcServer, &s)
	reflection.Register(grpcServer)

	log.Printf("starting gRPC server on %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server: %s", err)
	}
}

func withLogUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(loggingInterceptor)
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n", info.FullMethod, time.Since(start), err)
	return h, err
}
