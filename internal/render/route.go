//revive:disable:package-comments
package render

import (
	"fmt"
	"slices"
)

// Files maps file paths to their contents.
type Files map[string][]byte

// Func renders source files into a YAML manifest blob.
type Func func(Files) ([]byte, error)

// MatchFunc reports whether the source files match a rendering backend.
type MatchFunc func(Files) bool

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

func match(files Files) func(entry) bool {
	return func(e entry) bool { return e.match(files) }
}

// Render dispatches to the first backend whose match predicate returns true.
func (r *Router) Render(files Files) ([]byte, error) {
	idx := slices.IndexFunc(r.backends, match(files))

	if idx < 0 {
		return nil, fmt.Errorf("no matching renderer")
	}

	return r.backends[idx].fn(files)
}
