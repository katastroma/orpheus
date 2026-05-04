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
	log       *slog.Logger
	router    *render.Router
	streamFn  orderer.StreamFunc
	chunkSize int
}

// New creates a Service with the given logger, router, stream function, and chunk size.
func New(log *slog.Logger, router *render.Router, streamFn orderer.StreamFunc, chunkSize int) *Service {
	return &Service{log: log, router: router, streamFn: streamFn, chunkSize: chunkSize}
}
