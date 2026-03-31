package helm_test

import (
	"testing"

	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/render/helm"
)

func TestMatch(t *testing.T) {
	files := render.Files{"Chart.yaml": []byte("name: test")}
	if !helm.Match(files) {
		t.Fatal("expected match for Chart.yaml")
	}
}

func TestMatch_NoChart(t *testing.T) {
	files := render.Files{"deployment.yaml": []byte("kind: Deployment")}
	if helm.Match(files) {
		t.Fatal("expected no match without Chart.yaml")
	}
}

func TestMatch_Empty(t *testing.T) {
	if helm.Match(render.Files{}) {
		t.Fatal("expected no match for empty files")
	}
}
