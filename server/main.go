package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/nathanows/ues/echo"
	"github.com/nathanows/ues/internal/pkg/middleware"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	addr := os.Getenv("SERVER_ADDR")
	fmt.Println(addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to create listener - err=%s", err)
	}

	grpcServer, err := buildServer()
	if err != nil {
		log.Fatalf("failed to build gRPC server - err=%s", err)
	}

	registerEchoService(grpcServer)

	log.Printf("starting gRPC server - addr=%s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start gRPC server - err=%s", err)
	}
}

func buildServer() (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile("certs/server-cert.pem", "certs/server-key.pem")
	if err != nil {
		return nil, fmt.Errorf("unable to load TLS key - err=%s", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds), unaryMiddleware()}

	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)

	return grpcServer, nil
}

// Builds considated grpc.ServerOption from list of injected middleware. The
// core gRPC lib allows only a single unary interceptor. The grpc_middleware
// package used here allows middleware chaining. Interceptors are invoked
// in the order specified from top to bottom.
func unaryMiddleware() grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		middleware.AuthInterceptor,
		middleware.LoggingInterceptor,
	))
}

func registerEchoService(server *grpc.Server) {
	echoSvc := echo.Service{}
	echo.RegisterEchoServiceServer(server, &echoSvc)
}
