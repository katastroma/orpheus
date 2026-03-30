//revive:disable:package-comments
package kustomize

import (
	"slices"

	"github.com/katastroma/orpheus/internal/render"
)

// Type is the registration name for the Kustomize backend.
const Type = "kustomize"

// markerFiles are the filenames that indicate a kustomization root.
var markerFiles = []string{"kustomization.yaml", "kustomization.yml", "Kustomization"}

func matches(files render.Files) func(string) bool {
	return func(name string) bool {
		_, ok := files[name]
		return ok
	}
}

// Match reports whether the files contain a kustomization marker file.
func Match(files render.Files) bool {
	return slices.ContainsFunc(markerFiles, matches(files))
}
