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
