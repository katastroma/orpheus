//revive:disable:package-comments
package render

import (
	"fmt"
	"slices"

	"github.com/go-git/go-billy/v5"
)

// Func renders a filesystem into Kubernetes manifests.
type Func func(billy.Filesystem) ([][]byte, error)

// MatchFunc reports whether the filesystem matches a rendering backend.
type MatchFunc func(billy.Filesystem) bool

type entry struct {
	name  string
	match MatchFunc
	fn    Func
}

// Router dispatches rendering to the first registered backend whose match
// predicate returns true. Registration order determines priority.
type Router struct {
	backends []entry
}

// Register adds a rendering backend with a match predicate.
func (r *Router) Register(name string, match MatchFunc, fn Func) {
	r.backends = append(r.backends, entry{name: name, match: match, fn: fn})
}

// Render dispatches to the first backend whose match predicate returns true.
func (r *Router) Render(fs billy.Filesystem) ([][]byte, error) {
	idx := slices.IndexFunc(r.backends, func(e entry) bool {
		return e.match(fs)
	})

	if idx < 0 {
		return nil, fmt.Errorf("no matching renderer")
	}

	return r.backends[idx].fn(fs)
}
