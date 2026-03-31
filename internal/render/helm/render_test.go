package helm_test

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/render/helm"
)

func TestRender(t *testing.T) {
	files := render.Files{
		"Chart.yaml": []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"templates/configmap.yaml": []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-config
data:
  key: value
`),
	}

	out, err := helm.Render(files)
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
	files := render.Files{
		"Chart.yaml":  []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"values.yaml": []byte("replicas: 3\n"),
		"templates/deployment.yaml": []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
spec:
  replicas: {{ .Values.replicas }}
`),
	}

	out, err := helm.Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(out), "replicas: 3") {
		t.Errorf("expected values to be applied, got %q", string(out))
	}
}

func TestRender_InvalidChart(t *testing.T) {
	files := render.Files{
		"Chart.yaml": []byte("not valid chart"),
	}

	if _, err := helm.Render(files); err == nil {
		t.Fatal("expected error for invalid chart")
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	files := render.Files{
		"Chart.yaml": []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"templates/bad.yaml": []byte("{{ .Nonexistent.Deeply.Nested }}"),
	}

	if _, err := helm.Render(files); err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestRender_EmptyFiles(t *testing.T) {
	if _, err := helm.Render(render.Files{}); err == nil {
		t.Fatal("expected error for empty files")
	}
}
