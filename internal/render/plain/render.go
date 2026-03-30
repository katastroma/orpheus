//revive:disable:package-comments
package plain

import (
	"bytes"
	"path/filepath"

	"github.com/katastroma/orpheus/internal/render"
)

const yamlSeparator = "---"

// Render splits all YAML files into individual documents and returns each
// as a separate manifest.
func Render(files render.Files) ([][]byte, error) {
	var manifests [][]byte

	for path, content := range files {
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		for doc := range bytes.SplitSeq(content, []byte(yamlSeparator)) {
			trimmed := bytes.TrimSpace(doc)
			if len(trimmed) > 0 {
				manifests = append(manifests, trimmed)
			}
		}
	}

	return manifests, nil
}
