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
	renderFn render.Func
	streamFn orderer.StreamFunc
}

// New creates a Service with the given logger, render, and stream functions.
func New(log *slog.Logger, renderFn render.Func, streamFn orderer.StreamFunc) *Service {
	return &Service{log: log, renderFn: renderFn, streamFn: streamFn}
}
