//revive:disable:package-comments
package render

import (
	"fmt"
	"io"

	pb "github.com/katastroma/keleustes"
)

// Backend receives source content and renders it into a YAML manifest blob.
type Backend interface {
	// Receive buffers the stream into a format-specific representation.
	Receive(io.Reader) error
	// Render produces a YAML manifest blob from the received content.
	Render() ([]byte, error)
}

// Router dispatches rendering to the backend registered for a given type.
type Router struct {
	backends map[pb.RendererType]Backend
}

// Register adds a rendering backend for a renderer type.
func (r *Router) Register(rendererType pb.RendererType, backend Backend) {
	if r.backends == nil {
		r.backends = make(map[pb.RendererType]Backend)
	}
	r.backends[rendererType] = backend
}

// Lookup returns the backend registered for the given type.
func (r *Router) Lookup(rendererType pb.RendererType) (Backend, error) {
	backend, ok := r.backends[rendererType]
	if !ok {
		return nil, fmt.Errorf("no renderer registered for %s", rendererType)
	}

	return backend, nil
}
