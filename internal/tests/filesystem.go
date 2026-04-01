//revive:disable:package-comments
package tests

import (
	"fmt"
	"io"
)

// ErrReader is an io.Reader that always returns an error.
type ErrReader struct{}

// Read always returns a read error.
func (ErrReader) Read(_ []byte) (int, error) { return 0, fmt.Errorf("read failed") }

var _ io.Reader = ErrReader{}

// ErrWriter is an io.Writer that always returns an error.
type ErrWriter struct{}

// Write always returns a write error.
func (ErrWriter) Write(_ []byte) (int, error) { return 0, fmt.Errorf("write failed") }

var _ io.Writer = ErrWriter{}

// FailAfterNWriter is an io.Writer that succeeds for the first N bytes
// then returns an error.
type FailAfterNWriter struct {
	N       int
	written int
}

// Write succeeds until N bytes have been written, then returns an error.
func (w *FailAfterNWriter) Write(p []byte) (int, error) {
	if w.written+len(p) > w.N {
		return 0, fmt.Errorf("write failed after %d bytes", w.N)
	}
	w.written += len(p)
	return len(p), nil
}

var _ io.Writer = (*FailAfterNWriter)(nil)
