//revive:disable:package-comments
package plain

import "github.com/katastroma/orpheus/internal/render"

// Type is the registration name for the plain YAML backend.
const Type = "plain"

// Match always returns true. Plain is the catch-all backend — any set of
// YAML files is a valid render target.
func Match(_ render.Files) bool {
	return true
}
