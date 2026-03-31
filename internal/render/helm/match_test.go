package helm

import (
	"testing"

	"github.com/katastroma/orpheus/internal/render"
)

func TestMatch(t *testing.T) {
	files := render.Files{"Chart.yaml": []byte("name: test")}
	if !Match(files) {
		t.Fatal("expected match for Chart.yaml")
	}
}

func TestMatch_NoChart(t *testing.T) {
	files := render.Files{"deployment.yaml": []byte("kind: Deployment")}
	if Match(files) {
		t.Fatal("expected no match without Chart.yaml")
	}
}

func TestMatch_Empty(t *testing.T) {
	if Match(render.Files{}) {
		t.Fatal("expected no match for empty files")
	}
}
