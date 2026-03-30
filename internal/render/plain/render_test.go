package plain

import (
	"testing"

	"github.com/katastroma/orpheus/internal/render"
)

func TestRender(t *testing.T) {
	files := render.Files{"deployment.yaml": []byte("kind: Deployment")}

	manifests, err := Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(manifests))
	}

	if string(manifests[0]) != "kind: Deployment" {
		t.Errorf("expected %q, got %q", "kind: Deployment", string(manifests[0]))
	}
}

func TestRender_MultiDocument(t *testing.T) {
	files := render.Files{"manifests.yaml": []byte("kind: Namespace\n---\nkind: Deployment")}

	manifests, err := Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 2 {
		t.Fatalf("expected 2 manifests, got %d", len(manifests))
	}
}

func TestRender_SkipsEmptyDocuments(t *testing.T) {
	files := render.Files{"app.yaml": []byte("---\nkind: Service\n---\n\n---\n")}

	manifests, err := Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest (empty docs skipped), got %d", len(manifests))
	}
}

func TestRender_SkipsNonYAML(t *testing.T) {
	files := render.Files{
		"README.md":  []byte("# Hello"),
		"app.yaml":   []byte("kind: Pod"),
	}

	manifests, err := Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(manifests))
	}
}

func TestRender_YmlExtension(t *testing.T) {
	files := render.Files{"service.yml": []byte("kind: Service")}

	manifests, err := Render(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(manifests))
	}
}

func TestRender_EmptyFiles(t *testing.T) {
	manifests, err := Render(render.Files{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 0 {
		t.Errorf("expected 0 manifests, got %d", len(manifests))
	}
}
