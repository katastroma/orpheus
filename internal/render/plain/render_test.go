package plain

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/render"
)

func TestRender(t *testing.T) {
	files := render.Files{"deployment.yaml": []byte("kind: Deployment")}

	out, err := Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(out), "kind: Deployment") {
		t.Errorf("expected Deployment in output, got %q", string(out))
	}
}

func TestRender_MultipleFiles(t *testing.T) {
	files := render.Files{
		"a.yaml": []byte("kind: Namespace"),
		"b.yaml": []byte("kind: Deployment"),
	}

	out, err := Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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
	files := render.Files{
		"README.md": []byte("# Hello"),
		"app.yaml":  []byte("kind: Pod"),
	}

	out, err := Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(string(out), "Hello") {
		t.Error("expected non-YAML file to be skipped")
	}

	if !strings.Contains(string(out), "kind: Pod") {
		t.Error("expected YAML file in output")
	}
}

func TestRender_YmlExtension(t *testing.T) {
	files := render.Files{"service.yml": []byte("kind: Service")}

	out, err := Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(out), "kind: Service") {
		t.Errorf("expected Service in output, got %q", string(out))
	}
}

func TestRender_EmptyFiles(t *testing.T) {
	out, err := Render(render.Files{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out) != 0 {
		t.Errorf("expected empty output, got %q", string(out))
	}
}
