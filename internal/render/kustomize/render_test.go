package kustomize_test

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/render/kustomize"
)

func TestRender(t *testing.T) {
	files := render.Files{
		"kustomization.yaml": []byte("resources:\n- deployment.yaml\n"),
		"deployment.yaml": []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
`),
	}

	out, err := kustomize.Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(out), "kind: Deployment") {
		t.Errorf("expected Deployment in output, got %q", string(out))
	}
}

func TestRender_MultipleResources(t *testing.T) {
	files := render.Files{
		"kustomization.yaml": []byte("resources:\n- deployment.yaml\n- service.yaml\n"),
		"deployment.yaml": []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
`),
		"service.yaml": []byte(`apiVersion: v1
kind: Service
metadata:
  name: test
`),
	}

	out, err := kustomize.Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(out), "kind: Deployment") {
		t.Error("expected Deployment in output")
	}
	if !strings.Contains(string(out), "kind: Service") {
		t.Error("expected Service in output")
	}
}

func TestRender_InvalidKustomization(t *testing.T) {
	files := render.Files{
		"kustomization.yaml": []byte("not valid kustomization"),
	}

	if _, err := kustomize.Render(files); err == nil {
		t.Fatal("expected error for invalid kustomization")
	}
}

func TestRender_MissingResource(t *testing.T) {
	files := render.Files{
		"kustomization.yaml": []byte("resources:\n- nonexistent.yaml\n"),
	}

	if _, err := kustomize.Render(files); err == nil {
		t.Fatal("expected error for missing resource reference")
	}
}

func TestRender_InvalidFilePath(t *testing.T) {
	files := render.Files{
		"../invalid_file": []byte(""),
	}

	if _, err := kustomize.Render(files); err == nil {
		t.Fatal("expected error for invalid file path")
	}
}
