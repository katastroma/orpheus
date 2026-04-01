package kustomize_test

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/render/kustomize"
	"github.com/katastroma/orpheus/internal/tests"
)

func TestRender(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"kustomization.yaml": []byte("resources:\n- deployment.yaml\n"),
		"deployment.yaml": []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: test\n"),
	})

	out, err := kustomize.Render(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(out), "kind: Deployment") {
		t.Errorf("expected Deployment in output, got %q", string(out))
	}
}

func TestRender_MultipleResources(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"kustomization.yaml": []byte("resources:\n- deployment.yaml\n- service.yaml\n"),
		"deployment.yaml":    []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: test\n"),
		"service.yaml":       []byte("apiVersion: v1\nkind: Service\nmetadata:\n  name: test\n"),
	})

	out, err := kustomize.Render(r)
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
	r := tests.TarReader(t, map[string][]byte{
		"kustomization.yaml": []byte("not valid kustomization"),
	})

	if _, err := kustomize.Render(r); err == nil {
		t.Fatal("expected error for invalid kustomization")
	}
}

func TestRender_CorruptArchive(t *testing.T) {
	if _, err := kustomize.Render(strings.NewReader("not a tar")); err == nil {
		t.Fatal("expected error for corrupt archive")
	}
}

func TestRender_MissingResource(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"kustomization.yaml": []byte("resources:\n- nonexistent.yaml\n"),
	})

	if _, err := kustomize.Render(r); err == nil {
		t.Fatal("expected error for missing resource reference")
	}
}
