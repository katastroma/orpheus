package kustomize

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/tests"
)

func TestExtractToFS_ReadError(t *testing.T) {
	if _, err := extractToFS(tests.ErrReader{}); err == nil {
		t.Fatal("expected error when reader fails")
	}
}

func TestExtractToFS_TruncatedEntry(t *testing.T) {
	r := tests.TruncatedTarReader(t, "app.yaml", 1000)

	if _, err := extractToFS(r); err == nil {
		t.Fatal("expected error for truncated tar entry")
	}
}

func TestExtractToFS_CorruptArchive(t *testing.T) {
	if _, err := extractToFS(strings.NewReader("not a tar")); err == nil {
		t.Fatal("expected error for corrupt archive")
	}
}

func TestExtractToFS_SkipsDirectories(t *testing.T) {
	r := tests.TarReaderWithDir(t, "deploy/", map[string][]byte{
		"deploy/app.yaml": []byte("kind: Pod"),
	})

	fSys, err := extractToFS(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fSys.Exists("deploy/app.yaml") {
		t.Fatal("expected deploy/app.yaml in filesystem")
	}
}

func TestExtractToFS_CreatesSubdirectories(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"sub/dir/app.yaml": []byte("kind: Pod"),
	})

	fSys, err := extractToFS(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fSys.Exists("sub/dir/app.yaml") {
		t.Fatal("expected sub/dir/app.yaml in filesystem")
	}
}

func TestExtractToFS_InvalidPath(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"../escape.yaml": []byte("kind: Pod"),
	})

	if _, err := extractToFS(r); err == nil {
		t.Fatal("expected error for path traversal")
	}
}

func TestExtractToFS_WriteFileError(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"": []byte("kind: Pod"),
	})

	if _, err := extractToFS(r); err == nil {
		t.Fatal("expected error for empty file name")
	}
}

func TestExtractToFS_RootFile(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"app.yaml": []byte("kind: Pod"),
	})

	fSys, err := extractToFS(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fSys.Exists("app.yaml") {
		t.Fatal("expected app.yaml in filesystem")
	}
}
