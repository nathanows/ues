package main

import (
	"log"
	"net"

	"github.com/nathanows/ues/echo"
	"github.com/nathanows/ues/internal/pkg/middleware"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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

	s := echo.Service{}

	creds, err := credentials.NewServerTLSFromFile("certs/server-cert.pem", "certs/server-key.pem")
	if err != nil {
		log.Fatalf("unable to load TLS key - err=%s", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds), unaryMiddleware()}

	grpcServer := grpc.NewServer(opts...)
	echo.RegisterEchoServiceServer(grpcServer, &s)
	reflection.Register(grpcServer)

	log.Printf("starting gRPC server - port=%s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server - err=%s", err)
	}
}

func unaryMiddleware() grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		middleware.AuthInterceptor,
		middleware.LoggingInterceptor,
	))
}
