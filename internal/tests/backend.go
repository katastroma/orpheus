//revive:disable:package-comments
package tests

import (
	"io"

	"github.com/katastroma/orpheus/internal/render"
)

var _ render.Backend = (*MockBackend)(nil)

// MockBackend implements render.Backend for testing.
type MockBackend struct {
	// ReceiveErr is returned by Receive when set.
	ReceiveErr error
	// RenderResult is returned by Render on success.
	RenderResult []byte
	// RenderErr is returned by Render when set.
	RenderErr error
}

// Receive returns the configured error.
func (b *MockBackend) Receive(_ io.Reader) error {
	return b.ReceiveErr
}

// Render returns the configured result or error.
func (b *MockBackend) Render() ([]byte, error) {
	if b.RenderErr != nil {
		return nil, b.RenderErr
	}
	return b.RenderResult, nil
}
