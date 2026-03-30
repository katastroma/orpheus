package extract

import (
	"archive/tar"
	"bytes"
	"testing"

	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/tests"
)

func writeTarEntry(t *testing.T, tw *tar.Writer, name, content string) {
	t.Helper()

	if err := tw.WriteHeader(&tar.Header{
		Name:     name,
		Size:     int64(len(content)),
		Typeflag: tar.TypeReg,
	}); err != nil {
		t.Fatalf("writing tar header for %s: %v", name, err)
	}

	if _, err := tw.Write([]byte(content)); err != nil {
		t.Fatalf("writing tar content for %s: %v", name, err)
	}
}

func buildTar(t *testing.T, files render.Files) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for name, content := range files {
		writeTarEntry(t, tw, name, string(content))
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar writer: %v", err)
	}

	return &buf
}

func TestTar(t *testing.T) {
	buf := buildTar(t, render.Files{"values.yaml": []byte("key: value")})

	files, err := Tar(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, ok := files["values.yaml"]
	if !ok {
		t.Fatal("expected values.yaml in extracted files")
	}

	if string(content) != "key: value" {
		t.Errorf("expected %q, got %q", "key: value", string(content))
	}
}

func TestTar_MultipleFiles(t *testing.T) {
	buf := buildTar(t, render.Files{
		"a.yaml": []byte("a"),
		"b.yaml": []byte("b"),
	})

	files, err := Tar(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestTar_SkipsDirectories(t *testing.T) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	if err := tw.WriteHeader(&tar.Header{
		Name:     "deploy/",
		Typeflag: tar.TypeDir,
	}); err != nil {
		t.Fatalf("writing dir header: %v", err)
	}

	writeTarEntry(t, tw, "deploy/app.yaml", "kind: Deployment")

	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar: %v", err)
	}

	files, err := Tar(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := files["deploy/"]; ok {
		t.Error("expected directory entry to be skipped")
	}

	if _, ok := files["deploy/app.yaml"]; !ok {
		t.Error("expected deploy/app.yaml in extracted files")
	}
}

func TestTar_EmptyArchive(t *testing.T) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar: %v", err)
	}

	files, err := Tar(&buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestTar_TruncatedEntry(t *testing.T) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	// Write a header claiming 1000 bytes but don't write the content.
	if err := tw.WriteHeader(&tar.Header{
		Name:     "big.yaml",
		Size:     1000,
		Typeflag: tar.TypeReg,
	}); err != nil {
		t.Fatalf("writing header: %v", err)
	}

	// Truncate: close the underlying buffer without writing content or
	// closing the tar writer, producing an incomplete entry.

	if _, err := Tar(&buf); err == nil {
		t.Fatal("expected error for truncated tar entry")
	}
}

func TestTar_CorruptArchive(t *testing.T) {
	buf := bytes.NewBufferString("not a tar archive")

	if _, err := Tar(buf); err == nil {
		t.Fatal("expected error for corrupt tar")
	}
}

func TestTar_ReadError(t *testing.T) {
	if _, err := Tar(tests.ErrReader{}); err == nil {
		t.Fatal("expected error when reader fails")
	}
}
