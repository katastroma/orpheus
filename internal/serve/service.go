//revive:disable:package-comments
package serve

import (
	pb "github.com/katastroma/keleustes"

	"github.com/katastroma/orpheus/internal/orderer"
	"github.com/katastroma/orpheus/internal/render"
)

// Service implements the keleustes RendererServiceServer.
type Service struct {
	pb.UnimplementedRendererServiceServer
	renderFn render.Func
	streamFn orderer.StreamFunc
}

// New creates a Service with the given render and stream functions.
func New(renderFn render.Func, streamFn orderer.StreamFunc) *Service {
	return &Service{renderFn: renderFn, streamFn: streamFn}
}
