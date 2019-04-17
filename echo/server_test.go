package echo_test

import (
	context "context"
	"testing"

	"github.com/nathanows/ues/echo"
)

func TestEcho(t *testing.T) {
	s := &echo.EchoService{}

	tests := []struct {
		req  string
		resp string
	}{
		{req: "hello", resp: "hello"},
		{req: "123", resp: "123"},
	}

	for _, tt := range tests {
		req := &echo.EchoRequest{Message: tt.req}
		resp, err := s.Echo(context.Background(), req)
		if err != nil {
			t.Errorf("Echo({Message: %s}) got unexpected error: %s", tt.req, err)
		}
		if resp.Message != tt.resp {
			t.Errorf("Echo({Message: %s})=%s, expected %s", tt.req, resp.Message, tt.resp)
		}
	}
}
