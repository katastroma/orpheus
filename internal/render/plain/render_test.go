package plain

import (
	"fmt"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"

	"github.com/katastroma/orpheus/internal/tests"
)

func createTestFile(t *testing.T, fs billy.Filesystem, path, content string) {
	t.Helper()

	f, err := fs.Create(path)
	if err != nil {
		t.Fatalf("creating %s: %v", path, err)
	}

	if _, err = f.Write([]byte(content)); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}

	if err = f.Close(); err != nil {
		t.Fatalf("closing %s: %v", path, err)
	}
}

func TestRender(t *testing.T) {
	fs := memfs.New()
	createTestFile(t, fs, "deployment.yaml", "kind: Deployment")

	manifests, err := Render(fs)
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
	fs := memfs.New()
	createTestFile(t, fs, "manifests.yaml", "kind: Namespace\n---\nkind: Deployment")

	manifests, err := Render(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 2 {
		t.Fatalf("expected 2 manifests, got %d", len(manifests))
	}

	if string(manifests[0]) != "kind: Namespace" {
		t.Errorf("expected %q, got %q", "kind: Namespace", string(manifests[0]))
	}

	if string(manifests[1]) != "kind: Deployment" {
		t.Errorf("expected %q, got %q", "kind: Deployment", string(manifests[1]))
	}
}

func TestRender_SkipsEmptyDocuments(t *testing.T) {
	fs := memfs.New()
	createTestFile(t, fs, "app.yaml", "---\nkind: Service\n---\n\n---\n")

	manifests, err := Render(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest (empty docs skipped), got %d", len(manifests))
	}
}

func TestRender_SkipsNonYAML(t *testing.T) {
	fs := memfs.New()
	createTestFile(t, fs, "README.md", "# Hello")
	createTestFile(t, fs, "app.yaml", "kind: Pod")

	manifests, err := Render(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(manifests))
	}
}

func TestRender_YmlExtension(t *testing.T) {
	fs := memfs.New()
	createTestFile(t, fs, "service.yml", "kind: Service")

	manifests, err := Render(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(manifests))
	}
}

func TestRender_EmptyFilesystem(t *testing.T) {
	fs := memfs.New()

	manifests, err := Render(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 0 {
		t.Errorf("expected 0 manifests, got %d", len(manifests))
	}
}

func TestWalker_Visit_WalkError(t *testing.T) {
	w := &walker{fs: memfs.New()}

	walkErr := fmt.Errorf("walk error")
	if err := w.visit("path", nil, walkErr); err != walkErr {
		t.Fatalf("expected walk error, got %v", err)
	}
}

func TestRender_OpenError(t *testing.T) {
	backing := memfs.New()
	createTestFile(t, backing, "app.yaml", "kind: Pod")
	fs := &tests.OpenErrorFS{Filesystem: backing}

	if _, err := Render(fs); err == nil {
		t.Fatal("expected error when file open fails")
	}
}

func TestRender_ReadError(t *testing.T) {
	backing := memfs.New()
	createTestFile(t, backing, "app.yaml", "kind: Pod")
	fs := &tests.ErrorFS{Filesystem: backing}

	if _, err := Render(fs); err == nil {
		t.Fatal("expected error when file read fails")
	}
}
