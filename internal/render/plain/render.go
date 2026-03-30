//revive:disable:package-comments
package plain

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
)

const yamlSeparator = "---"

// Render reads all .yaml and .yml files from fs and returns each YAML
// document as a separate manifest.
func Render(fs billy.Filesystem) ([][]byte, error) {
	w := &walker{fs: fs}

	if err := util.Walk(fs, ".", w.visit); err != nil {
		return nil, err
	}

	return w.manifests, nil
}

// walker collects manifests from YAML files during a filesystem walk.
type walker struct {
	fs        billy.Filesystem
	manifests [][]byte
}

// visit processes a single walk entry, reading YAML documents from files.
func (w *walker) visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
		return nil
	}

	docs, readErr := readDocuments(w.fs, path)
	if readErr != nil {
		return readErr
	}

	w.manifests = append(w.manifests, docs...)
	return nil
}

// readDocuments reads a YAML file and splits it into individual documents.
func readDocuments(fs billy.Filesystem, path string) ([][]byte, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var docs [][]byte
	for doc := range bytes.SplitSeq(content, []byte(yamlSeparator)) {
		trimmed := bytes.TrimSpace(doc)
		if len(trimmed) > 0 {
			docs = append(docs, trimmed)
		}
	}

	return docs, nil
}
