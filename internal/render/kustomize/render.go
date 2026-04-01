//revive:disable:package-comments
package kustomize

import (
	"archive/tar"
	"fmt"
	"io"
	"path/filepath"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// Backend implements render.Backend for Kustomize sources.
type Backend struct {
	fSys filesys.FileSystem
}

// New creates a Kustomize rendering backend.
func New() *Backend {
	return &Backend{}
}

// Receive extracts the tar stream into an in-memory filesystem.
func (b *Backend) Receive(r io.Reader) error {
	fSys, err := extractToFS(r)
	if err != nil {
		return err
	}

	b.fSys = fSys
	return nil
}

// Render runs kustomize build on the received filesystem.
func (b *Backend) Render() ([]byte, error) {
	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	resMap, err := k.Run(b.fSys, ".")
	if err != nil {
		return nil, fmt.Errorf("running kustomize: %w", err)
	}

	return resMap.AsYaml()
}

func extractToFS(r io.Reader) (filesys.FileSystem, error) {
	fSys := filesys.MakeFsInMemory()
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

		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", header.Name, err)
		}

		dir := filepath.Dir(header.Name)
		if dir != "." {
			if err = fSys.MkdirAll(dir); err != nil {
				return nil, fmt.Errorf("creating directory %s: %w", dir, err)
			}
		}

		if err = fSys.WriteFile(header.Name, data); err != nil {
			return nil, fmt.Errorf("writing %s: %w", header.Name, err)
		}
	}

	return fSys, nil
}
