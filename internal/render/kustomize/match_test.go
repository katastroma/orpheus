package kustomize_test

import (
	"testing"

	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/render/kustomize"
)

func TestMatch_KustomizationYaml(t *testing.T) {
	files := render.Files{"kustomization.yaml": []byte("resources:")}
	if !kustomize.Match(files) {
		t.Fatal("expected match for kustomization.yaml")
	}
}

func TestMatch_KustomizationYml(t *testing.T) {
	files := render.Files{"kustomization.yml": []byte("resources:")}
	if !kustomize.Match(files) {
		t.Fatal("expected match for kustomization.yml")
	}
}

func TestMatch_KustomizationCapitalized(t *testing.T) {
	files := render.Files{"Kustomization": []byte("resources:")}
	if !kustomize.Match(files) {
		t.Fatal("expected match for Kustomization")
	}
}

func TestMatch_NoMarker(t *testing.T) {
	files := render.Files{"deployment.yaml": []byte("kind: Deployment")}
	if kustomize.Match(files) {
		t.Fatal("expected no match without kustomization marker")
	}
}

func TestMatch_Empty(t *testing.T) {
	if kustomize.Match(render.Files{}) {
		t.Fatal("expected no match for empty files")
	}
}
