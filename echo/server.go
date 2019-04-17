package echo

import context "context"

// EchoService implements gRPC echo.
type EchoService struct {
}

// Echo responds with a message body matching that of the request message body
func (s *EchoService) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{Message: req.Message}, nil
}
