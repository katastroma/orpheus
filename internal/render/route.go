//revive:disable:package-comments
package render

import (
	"fmt"
	"io"

	pb "github.com/katastroma/keleustes"
)

// Func renders source content into a YAML manifest blob.
type Func func(io.Reader) ([]byte, error)

// DispatchFunc dispatches rendering to a backend based on renderer type.
type DispatchFunc func(pb.RendererType, io.Reader) ([]byte, error)

// Router dispatches rendering to the backend registered for a given type.
type Router struct {
	backends map[pb.RendererType]Func
}

// Register adds a rendering backend for a renderer type.
func (r *Router) Register(rendererType pb.RendererType, fn Func) {
	if r.backends == nil {
		r.backends = make(map[pb.RendererType]Func)
	}
	r.backends[rendererType] = fn
}

// Render dispatches to the backend registered for the given type.
func (r *Router) Render(rendererType pb.RendererType, reader io.Reader) ([]byte, error) {
	fn, ok := r.backends[rendererType]
	if !ok {
		return nil, fmt.Errorf("no renderer registered for %s", rendererType)
	}

	return fn(reader)
}
