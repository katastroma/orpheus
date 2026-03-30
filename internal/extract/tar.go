//revive:disable:package-comments
package extract

import (
	"archive/tar"
	"fmt"
	"io"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

// Tar reads a tar archive from r and writes its entries into fs.
func Tar(r io.Reader, fs billy.Filesystem) error {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading tar header: %w", err)
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		dir := filepath.Dir(header.Name)
		if dir != "." {
			if err = fs.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("creating directory %s: %w", dir, err)
			}
		}

		f, err := fs.Create(header.Name)
		if err != nil {
			return fmt.Errorf("creating file %s: %w", header.Name, err)
		}

		if _, err = io.Copy(f, tr); err != nil {
			f.Close()
			return fmt.Errorf("writing file %s: %w", header.Name, err)
		}

		if err = f.Close(); err != nil {
			return fmt.Errorf("closing file %s: %w", header.Name, err)
		}
	}
}
