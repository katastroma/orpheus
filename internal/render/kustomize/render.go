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

// Render extracts the tar stream into an in-memory filesystem and runs
// kustomize build on it.
func Render(r io.Reader) ([]byte, error) {
	fSys, err := extractToFS(r)
	if err != nil {
		return nil, err
	}

	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	resMap, err := k.Run(fSys, ".")
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
