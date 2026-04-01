package plain_test

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/render/plain"
	"github.com/katastroma/orpheus/internal/tests"
)

func TestRender(t *testing.T) {
	b := plain.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{
		"deployment.yaml": []byte("kind: Deployment"),
	})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	out, err := b.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if !strings.Contains(string(out), "kind: Deployment") {
		t.Errorf("expected Deployment in output, got %q", string(out))
	}
}

func TestRender_MultipleFiles(t *testing.T) {
	b := plain.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{
		"a.yaml": []byte("kind: Namespace"),
		"b.yaml": []byte("kind: Deployment"),
	})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	out, err := b.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if !strings.Contains(string(out), "kind: Namespace") {
		t.Error("expected Namespace in output")
	}
	if !strings.Contains(string(out), "kind: Deployment") {
		t.Error("expected Deployment in output")
	}
	if !strings.Contains(string(out), "---") {
		t.Error("expected document separator between files")
	}
}

func TestRender_SkipsNonYAML(t *testing.T) {
	b := plain.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{
		"README.md": []byte("# Hello"),
		"app.yaml":  []byte("kind: Pod"),
	})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	out, err := b.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if strings.Contains(string(out), "Hello") {
		t.Error("expected non-YAML file to be skipped")
	}
	if !strings.Contains(string(out), "kind: Pod") {
		t.Error("expected YAML file in output")
	}
}

func TestRender_YmlExtension(t *testing.T) {
	b := plain.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{
		"service.yml": []byte("kind: Service"),
	})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	out, err := b.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if !strings.Contains(string(out), "kind: Service") {
		t.Errorf("expected Service in output, got %q", string(out))
	}
}

func TestRender_EmptyArchive(t *testing.T) {
	b := plain.New()

	if err := b.Receive(tests.TarReader(t, map[string][]byte{})); err != nil {
		t.Fatalf("receive: %v", err)
	}

	out, err := b.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if len(out) != 0 {
		t.Errorf("expected empty output, got %q", string(out))
	}
}
