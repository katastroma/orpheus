//revive:disable:package-comments
package plain

import (
	"archive/tar"
	"fmt"
	"io"
	"path/filepath"
)

// Backend implements render.Backend for plain YAML sources.
type Backend struct {
	manifests []byte
}

// New creates a plain YAML rendering backend.
func New() *Backend {
	return &Backend{}
}

// Receive extracts the tar stream and collects all YAML files.
func (b *Backend) Receive(r io.Reader) error {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading tar: %w", err)
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		ext := filepath.Ext(header.Name)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		content, err := io.ReadAll(tr)
		if err != nil {
			return fmt.Errorf("reading %s: %w", header.Name, err)
		}

		if len(b.manifests) > 0 {
			b.manifests = append(b.manifests, []byte("\n---\n")...)
		}
		b.manifests = append(b.manifests, content...)
	}
}

// Render returns the collected YAML manifests.
func (b *Backend) Render() ([]byte, error) {
	return b.manifests, nil
}
