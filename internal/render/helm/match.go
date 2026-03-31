//revive:disable:package-comments
package helm

import "github.com/katastroma/orpheus/internal/render"

// Type is the registration name for the Helm backend
const Type = "helm"

// Match reports whether the files contain a Chart.yaml
func Match(files render.Files) bool {
	_, ok := files["Chart.yaml"]
	return ok
}
