//revive:disable:package-comments
package plain

import (
	"archive/tar"
	"fmt"
	"io"
	"path/filepath"
)

// Render extracts the tar stream and concatenates all YAML files into a
// single manifest blob.
func Render(r io.Reader) ([]byte, error) {
	var out []byte
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading tar: %w", err)
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
			return nil, fmt.Errorf("reading %s: %w", header.Name, err)
		}

		if len(out) > 0 {
			out = append(out, []byte("\n---\n")...)
		}
		out = append(out, content...)
	}

	return out, nil
}
