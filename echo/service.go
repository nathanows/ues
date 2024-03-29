package echo

import "context"

// Service implements EchoServiceServer
type Service struct {
}

// Echo responds with a message body matching that of the request message body
func (s *Service) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{Message: req.Message}, nil
}
