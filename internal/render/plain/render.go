//revive:disable:package-comments
package plain

import (
	"path/filepath"

	"github.com/katastroma/orpheus/internal/render"
)

// Render concatenates all YAML files into a single manifest blob.
func Render(files render.Files) ([]byte, error) {
	var out []byte

	for path, content := range files {
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		if len(out) > 0 {
			out = append(out, []byte("\n---\n")...)
		}
		out = append(out, content...)
	}

	return out, nil
}
