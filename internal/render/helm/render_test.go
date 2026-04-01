package helm_test

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/render/helm"
	"github.com/katastroma/orpheus/internal/tests"
)

func TestRender(t *testing.T) {
	b := helm.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{
		"Chart.yaml":              []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"templates/configmap.yaml": []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: {{ .Release.Name }}-config\ndata:\n  key: value\n"),
	})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	out, err := b.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if !strings.Contains(string(out), "kind: ConfigMap") {
		t.Errorf("expected ConfigMap in output, got %q", string(out))
	}

	if !strings.Contains(string(out), "release-config") {
		t.Errorf("expected release name in output, got %q", string(out))
	}
}

func TestRender_WithValues(t *testing.T) {
	b := helm.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{
		"Chart.yaml":              []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"values.yaml":             []byte("replicas: 3\n"),
		"templates/deployment.yaml": []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: test\nspec:\n  replicas: {{ .Values.replicas }}\n"),
	})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	out, err := b.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if !strings.Contains(string(out), "replicas: 3") {
		t.Errorf("expected values to be applied, got %q", string(out))
	}
}

func TestRender_NestedChart(t *testing.T) {
	b := helm.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{
		"test/Chart.yaml":              []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"test/templates/configmap.yaml": []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\n"),
	})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	if _, err := b.Render(); err == nil {
		t.Fatal("expected error for chart not at archive root")
	}
}

func TestRender_InvalidChart(t *testing.T) {
	b := helm.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{
		"Chart.yaml": []byte("not valid chart"),
	})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	if _, err := b.Render(); err == nil {
		t.Fatal("expected error for invalid chart")
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	b := helm.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{
		"Chart.yaml":         []byte("apiVersion: v2\nname: test\nversion: 0.1.0\n"),
		"templates/bad.yaml": []byte("{{ .Nonexistent.Deeply.Nested }}"),
	})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	if _, err := b.Render(); err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestReceive_ReadError(t *testing.T) {
	b := helm.New()

	if err := b.Receive(tests.ErrReader{}); err == nil {
		t.Fatal("expected error when reader fails")
	}
}

func TestReceive_EmptyArchive(t *testing.T) {
	b := helm.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	if _, err := b.Render(); err == nil {
		t.Fatal("expected error for empty archive")
	}
}
