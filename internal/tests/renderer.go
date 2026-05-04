//revive:disable:package-comments
package tests

import (
	"context"
	"io"

	pb "github.com/katastroma/keleustes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MockRenderStreamServer implements pb.RendererService_RenderStreamServer for testing.
type MockRenderStreamServer struct {
	// Requests are returned sequentially by Recv.
	Requests []*pb.RenderStreamRequest
	// Sent collects messages passed to Send.
	Sent []*pb.RenderStreamResponse
	// SendErr is returned by Send when set.
	SendErr error
	// RecvErr is returned by Recv when set.
	RecvErr error
	// Ctx is returned by Context.
	Ctx     context.Context
	recvIdx int
	grpc.ServerStream
}

// Recv returns the next pre-loaded request, the configured error, or io.EOF.
func (s *MockRenderStreamServer) Recv() (*pb.RenderStreamRequest, error) {
	if s.RecvErr != nil {
		return nil, s.RecvErr
	}
	if s.recvIdx >= len(s.Requests) {
		return nil, io.EOF
	}
	req := s.Requests[s.recvIdx]
	s.recvIdx++
	return req, nil
}

// Send records a sent response or returns the configured error.
func (s *MockRenderStreamServer) Send(resp *pb.RenderStreamResponse) error {
	if s.SendErr != nil {
		return s.SendErr
	}
	s.Sent = append(s.Sent, resp)
	return nil
}

// Context returns the configured context.
func (s *MockRenderStreamServer) Context() context.Context {
	return s.Ctx
}

// SetHeader is a no-op.
func (s *MockRenderStreamServer) SetHeader(_ metadata.MD) error { return nil }

// SendHeader is a no-op.
func (s *MockRenderStreamServer) SendHeader(_ metadata.MD) error { return nil }

// SetTrailer is a no-op.
func (s *MockRenderStreamServer) SetTrailer(_ metadata.MD) {}

// SendMsg is a no-op.
func (s *MockRenderStreamServer) SendMsg(_ any) error { return nil }

// RecvMsg signals end of stream.
func (s *MockRenderStreamServer) RecvMsg(_ any) error { return io.EOF }

// MockRenderServer implements pb.RendererService_RenderServer for testing.
type MockRenderServer struct {
	// Requests are returned sequentially by Recv.
	Requests []*pb.RenderRequest
	// Responses collects messages passed to Send.
	Responses []*pb.RenderResponse
	// SendErr is returned by Send when set.
	SendErr error
	// RecvErr is returned by Recv when set.
	RecvErr error
	// Ctx is returned by Context.
	Ctx     context.Context
	recvIdx int
	grpc.ServerStream
}

// Recv returns the next pre-loaded request, the configured error, or io.EOF.
func (s *MockRenderServer) Recv() (*pb.RenderRequest, error) {
	if s.RecvErr != nil {
		return nil, s.RecvErr
	}
	if s.recvIdx >= len(s.Requests) {
		return nil, io.EOF
	}
	req := s.Requests[s.recvIdx]
	s.recvIdx++
	return req, nil
}

// SendAndClose records the response or returns the configured error.
func (s *MockRenderServer) SendAndClose(resp *pb.RenderResponse) error {
	if s.SendErr != nil {
		return s.SendErr
	}
	s.Responses = append(s.Responses, resp)
	return nil
}

// Context returns the configured context.
func (s *MockRenderServer) Context() context.Context {
	return s.Ctx
}

// SetHeader is a no-op.
func (s *MockRenderServer) SetHeader(_ metadata.MD) error { return nil }

// SendHeader is a no-op.
func (s *MockRenderServer) SendHeader(_ metadata.MD) error { return nil }

// SetTrailer is a no-op.
func (s *MockRenderServer) SetTrailer(_ metadata.MD) {}

// SendMsg is a no-op.
func (s *MockRenderServer) SendMsg(_ any) error { return nil }

// RecvMsg signals end of stream.
func (s *MockRenderServer) RecvMsg(_ any) error { return io.EOF }
