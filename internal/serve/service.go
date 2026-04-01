//revive:disable:package-comments
package serve

import (
	"log/slog"

	pb "github.com/katastroma/keleustes"

	"github.com/katastroma/orpheus/internal/orderer"
	"github.com/katastroma/orpheus/internal/render"
)

// Service implements the keleustes RendererServiceServer.
type Service struct {
	pb.UnimplementedRendererServiceServer
	log      *slog.Logger
	router   *render.Router
	streamFn orderer.StreamFunc
}

// New creates a Service with the given logger, router, and stream function.
func New(log *slog.Logger, router *render.Router, streamFn orderer.StreamFunc) *Service {
	return &Service{log: log, router: router, streamFn: streamFn}
}
