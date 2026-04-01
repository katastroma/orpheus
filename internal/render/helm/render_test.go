package helm_test

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/render/helm"
	"github.com/katastroma/orpheus/internal/tests"
)

func TestRender(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"Chart.yaml":              []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"templates/configmap.yaml": []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: {{ .Release.Name }}-config\ndata:\n  key: value\n"),
	})

	out, err := helm.Render(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(out), "kind: ConfigMap") {
		t.Errorf("expected ConfigMap in output, got %q", string(out))
	}

	if !strings.Contains(string(out), "release-config") {
		t.Errorf("expected release name in output, got %q", string(out))
	}
}

func TestRender_WithValues(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"Chart.yaml":              []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"values.yaml":             []byte("replicas: 3\n"),
		"templates/deployment.yaml": []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: test\nspec:\n  replicas: {{ .Values.replicas }}\n"),
	})

	out, err := helm.Render(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(out), "replicas: 3") {
		t.Errorf("expected values to be applied, got %q", string(out))
	}
}

func TestRender_NestedChart(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"test/Chart.yaml":              []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"test/templates/configmap.yaml": []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\n"),
	})

	if _, err := helm.Render(r); err == nil {
		t.Fatal("expected error for chart not at archive root")
	}
}

func TestRender_InvalidChart(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"Chart.yaml": []byte("not valid chart"),
	})

	if _, err := helm.Render(r); err == nil {
		t.Fatal("expected error for invalid chart")
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"Chart.yaml":         []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"templates/bad.yaml": []byte("{{ .Nonexistent.Deeply.Nested }}"),
	})

	if _, err := helm.Render(r); err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestRender_ReadError(t *testing.T) {
	if _, err := helm.Render(tests.ErrReader{}); err == nil {
		t.Fatal("expected error when reader fails")
	}
}

func TestRender_EmptyArchive(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{})

	if _, err := helm.Render(r); err == nil {
		t.Fatal("expected error for empty archive")
	}
}
