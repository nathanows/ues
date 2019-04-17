package middleware_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/nathanows/ues/internal/pkg/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func TestLoggingInternceptor(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	methodCall := "some.ServiceName/method"
	unaryInfo := &grpc.UnaryServerInfo{
		FullMethod: methodCall,
	}

	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return true, nil
	}

	ctx := context.Background()

	resp, err := middleware.LoggingInterceptor(ctx, "abc", unaryInfo, unaryHandler)
	if err != nil {
		t.Fatalf("unexpected failure")
	}

	if resp.(bool) != true {
		t.Fatalf("unexpected response: expected: true, got: %v", resp)
	}

	logMsg := buf.String()

	if !strings.Contains(logMsg, fmt.Sprintf("method=%s", methodCall)) {
		t.Fatalf("expected log message to contain method call info, got: %s", logMsg)
	}

	if !strings.Contains(logMsg, "duration=") {
		t.Fatalf("expected log message to contain duration tag, got: %s", logMsg)
	}

	if !strings.Contains(logMsg, "error=") {
		t.Fatalf("expected log message to contain error tag, got: %s", logMsg)
	}
}

func TestAuthInternceptor_Success(t *testing.T) {
	tokenKey := "TOKEN"
	currEnvToken := os.Getenv(tokenKey)
	testToken := "something"
	os.Setenv(tokenKey, testToken)
	defer func() {
		os.Setenv(tokenKey, currEnvToken)
	}()

	unaryInfo := &grpc.UnaryServerInfo{}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return true, nil
	}

	md := metadata.Pairs("authorization", testToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := middleware.AuthInterceptor(ctx, "abc", unaryInfo, unaryHandler)
	if err != nil {
		t.Errorf("unexpected failure: %#v", err)
	}

	if resp.(bool) != true {
		t.Errorf("unexpected response: expected: true, got: %v", resp)
	}
}

func TestAuthInternceptor_InvalidToken(t *testing.T) {
	tokenKey := "TOKEN"
	currEnvToken := os.Getenv(tokenKey)
	testToken := "incorrect-token"
	os.Setenv(tokenKey, "real-token")
	defer func() {
		os.Setenv(tokenKey, currEnvToken)
	}()

	unaryInfo := &grpc.UnaryServerInfo{}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return true, nil
	}

	md := metadata.Pairs("authorization", testToken)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := middleware.AuthInterceptor(ctx, "abc", unaryInfo, unaryHandler)

	if got, want := grpc.Code(err), codes.Unauthenticated; got != want {
		t.Errorf("expect grpc.Code to %s, but got %s", want, got)
	}

	if got := grpc.ErrorDesc(err); got != "invalid token" {
		t.Errorf("expected: %q, got: %q", "invalid token", got)
	}
}

func TestAuthInternceptor_NoToken(t *testing.T) {
	tokenKey := "TOKEN"
	currEnvToken := os.Getenv(tokenKey)
	os.Setenv(tokenKey, "real-token")
	defer func() {
		os.Setenv(tokenKey, currEnvToken)
	}()

	unaryInfo := &grpc.UnaryServerInfo{}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return true, nil
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})

	_, err := middleware.AuthInterceptor(ctx, "abc", unaryInfo, unaryHandler)

	if got, want := grpc.Code(err), codes.Unauthenticated; got != want {
		t.Errorf("expect grpc.Code to %s, but got %s", want, got)
	}

	if got := grpc.ErrorDesc(err); got != "invalid token" {
		t.Errorf("expected: %q, got: %q", "invalid token", got)
	}
}
