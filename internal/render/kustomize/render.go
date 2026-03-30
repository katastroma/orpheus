//revive:disable:package-comments
package kustomize

import (
	"fmt"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/katastroma/orpheus/internal/render"
)

// Render runs kustomize build on the given files and returns the rendered
// manifest blob.
func Render(files render.Files) ([]byte, error) {
	fSys, err := toKustomizeFS(files)
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

// toKustomizeFS populates a kustomize in-memory filesystem from the file map
func toKustomizeFS(files render.Files) (filesys.FileSystem, error) {
	fSys := filesys.MakeFsInMemory()

	for path, content := range files {
		if err := fSys.WriteFile(path, content); err != nil {
			return nil, fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return fSys, nil
}
