package plain_test

import (
	"testing"

	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/render/plain"
)

func TestMatch(t *testing.T) {
	if !plain.Match(render.Files{}) {
		t.Fatal("expected Match to return true")
	}
}
