//revive:disable:package-comments
package tests

import (
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
)

// ErrReader is an io.Reader that always returns an error.
type ErrReader struct{}

// Read always returns a read error.
func (ErrReader) Read(_ []byte) (int, error) { return 0, fmt.Errorf("read failed") }

var _ io.Reader = ErrReader{}

// OpenErrorFS wraps a real filesystem but errors on Open.
type OpenErrorFS struct{ billy.Filesystem }

// Open always returns an error.
func (fs *OpenErrorFS) Open(_ string) (billy.File, error) {
	return nil, fmt.Errorf("open denied")
}

// Stat delegates to the wrapped filesystem.
func (fs *OpenErrorFS) Stat(name string) (os.FileInfo, error) {
	return fs.Filesystem.Stat(name)
}

// ErrorFile implements billy.File with a Read that always errors.
type ErrorFile struct{ billy.File }

// Read always returns a disk error.
func (f *ErrorFile) Read(_ []byte) (int, error) { return 0, fmt.Errorf("disk error") }

// Close is a no-op.
func (f *ErrorFile) Close() error { return nil }

// ErrorFS wraps a real filesystem but returns an ErrorFile on Open.
type ErrorFS struct{ billy.Filesystem }

// Open returns an ErrorFile regardless of path.
func (fs *ErrorFS) Open(_ string) (billy.File, error) {
	return &ErrorFile{}, nil
}

// Stat delegates to the wrapped filesystem.
func (fs *ErrorFS) Stat(name string) (os.FileInfo, error) {
	return fs.Filesystem.Stat(name)
}

// MkdirErrorFS wraps a filesystem and errors on MkdirAll.
type MkdirErrorFS struct{ billy.Filesystem }

// MkdirAll always returns an error.
func (fs *MkdirErrorFS) MkdirAll(_ string, _ os.FileMode) error {
	return fmt.Errorf("mkdir denied")
}

// CreateErrorFS wraps a filesystem and errors on Create.
type CreateErrorFS struct{ billy.Filesystem }

// Create always returns an error.
func (fs *CreateErrorFS) Create(_ string) (billy.File, error) {
	return nil, fmt.Errorf("create denied")
}

// MkdirAll delegates to the wrapped filesystem.
func (fs *CreateErrorFS) MkdirAll(path string, perm os.FileMode) error {
	return fs.Filesystem.MkdirAll(path, perm)
}

// CloseErrorFile wraps a billy.File and errors on Close.
type CloseErrorFile struct{ billy.File }

// Close always returns an error.
func (f *CloseErrorFile) Close() error {
	return fmt.Errorf("close denied")
}

// CloseErrorFS wraps a filesystem and returns CloseErrorFile on Create.
type CloseErrorFS struct{ billy.Filesystem }

// Create returns a file that errors on Close.
func (fs *CloseErrorFS) Create(name string) (billy.File, error) {
	f, err := fs.Filesystem.Create(name)
	if err != nil {
		return nil, err
	}
	return &CloseErrorFile{File: f}, nil
}

// MkdirAll delegates to the wrapped filesystem.
func (fs *CloseErrorFS) MkdirAll(path string, perm os.FileMode) error {
	return fs.Filesystem.MkdirAll(path, perm)
}

// CopyErrorFile wraps a billy.File and errors on Write (used by io.Copy).
type CopyErrorFile struct{ billy.File }

// Write always returns an error.
func (f *CopyErrorFile) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("write denied")
}

// Close delegates to the wrapped file.
func (f *CopyErrorFile) Close() error {
	return f.File.Close()
}

// CopyErrorFS wraps a filesystem and returns CopyErrorFile on Create.
type CopyErrorFS struct{ billy.Filesystem }

// Create returns a file that errors on Write.
func (fs *CopyErrorFS) Create(name string) (billy.File, error) {
	f, err := fs.Filesystem.Create(name)
	if err != nil {
		return nil, err
	}
	return &CopyErrorFile{File: f}, nil
}

// MkdirAll delegates to the wrapped filesystem.
func (fs *CopyErrorFS) MkdirAll(path string, perm os.FileMode) error {
	return fs.Filesystem.MkdirAll(path, perm)
}
