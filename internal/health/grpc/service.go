//revive:disable:package-comments
package grpc

import (
	"context"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Service implements the gRPC health check protocol.
type Service struct {
	healthpb.UnimplementedHealthServer
}

// New returns a gRPC health service.
func New() *Service {
	return &Service{}
}

// Check reports service health.
func (s *Service) Check(
	_ context.Context,
	_ *healthpb.HealthCheckRequest,
) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
