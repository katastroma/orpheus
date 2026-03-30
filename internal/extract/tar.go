//revive:disable:package-comments
package extract

import (
	"archive/tar"
	"fmt"
	"io"

	"github.com/katastroma/orpheus/internal/render"
)

// Tar reads a tar archive from r and returns the extracted files.
func Tar(r io.Reader) (render.Files, error) {
	files := make(render.Files)
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			return files, nil
		}
		if err != nil {
			return nil, fmt.Errorf("reading tar header: %w", err)
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		content, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", header.Name, err)
		}

		files[header.Name] = content
	}
}
