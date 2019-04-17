// Package middleware provides grpc.UnaryServerInterceptor implementations
// providing intercepting hooks to be injected into the gRPC request lifecycle
// when instantiating new gRPC servers.
package middleware

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// LoggingInterceptor provides unary gRPC call middleware injecting standard
// semi-structured request logging. Incudes keypairs 'method', 'duration' and
// 'error' (<nil> if succesful request).
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	log.Printf("request processed - method=%s duration=%s error=%v\n", info.FullMethod, time.Since(start), err)
	return h, err
}

// AuthInterceptor provides unary gRPC call middleware which enforces the
// presence of a valid authorization token on all RPC's. Clients are to pass
// the preshared token via an "authorization" metadata key pair. Server-side
// the auth token is set using a TOKEN= environment variable.
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}
	if len(meta["authorization"]) != 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}
	if meta["authorization"][0] != os.Getenv("TOKEN") {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}
