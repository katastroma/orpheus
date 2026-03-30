package plain

import (
	"testing"

	"github.com/katastroma/orpheus/internal/render"
)

func TestMatch(t *testing.T) {
	if !Match(render.Files{}) {
		t.Fatal("expected Match to return true")
	}
}
