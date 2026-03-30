package extract

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"

	"github.com/katastroma/orpheus/internal/tests"
)

// writeTarEntry writes a single file entry to a tar writer.
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

// buildTar creates a tar archive from a map of path to content.
func buildTar(t *testing.T, files map[string]string) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for name, content := range files {
		writeTarEntry(t, tw, name, content)
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar writer: %v", err)
	}

	return &buf
}

func TestTar(t *testing.T) {
	buf := buildTar(t, map[string]string{
		"deploy/values.yaml": "key: value",
	})

	fs := memfs.New()
	if err := Tar(buf, fs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, err := fs.Open("deploy/values.yaml")
	if err != nil {
		t.Fatalf("opening extracted file: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("reading extracted file: %v", err)
	}

	if string(content) != "key: value" {
		t.Errorf("expected %q, got %q", "key: value", string(content))
	}
}

func TestTar_MultipleFiles(t *testing.T) {
	buf := buildTar(t, map[string]string{
		"a.yaml": "a",
		"b.yaml": "b",
	})

	fs := memfs.New()
	if err := Tar(buf, fs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, name := range []string{"a.yaml", "b.yaml"} {
		if _, err := fs.Stat(name); err != nil {
			t.Errorf("expected %s to exist: %v", name, err)
		}
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

	fs := memfs.New()
	if err := Tar(&buf, fs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, err := fs.Open("deploy/app.yaml")
	if err != nil {
		t.Fatalf("opening extracted file: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("reading extracted file: %v", err)
	}

	if string(content) != "kind: Deployment" {
		t.Errorf("expected %q, got %q", "kind: Deployment", string(content))
	}
}

func TestTar_EmptyArchive(t *testing.T) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar: %v", err)
	}

	fs := memfs.New()
	if err := Tar(&buf, fs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTar_CorruptArchive(t *testing.T) {
	buf := bytes.NewBufferString("not a tar archive")

	fs := memfs.New()
	if err := Tar(buf, fs); err == nil {
		t.Fatal("expected error for corrupt tar")
	}
}

func TestTar_RootLevelFile(t *testing.T) {
	buf := buildTar(t, map[string]string{
		"values.yaml": "key: value",
	})

	fs := memfs.New()
	if err := Tar(buf, fs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := fs.Stat("values.yaml"); err != nil {
		t.Fatalf("expected root-level file to exist: %v", err)
	}
}

func TestTar_MkdirError(t *testing.T) {
	buf := buildTar(t, map[string]string{
		"deploy/values.yaml": "key: value",
	})

	fs := &tests.MkdirErrorFS{Filesystem: memfs.New()}
	if err := Tar(buf, fs); err == nil {
		t.Fatal("expected error when mkdir fails")
	}
}

func TestTar_CreateError(t *testing.T) {
	buf := buildTar(t, map[string]string{
		"values.yaml": "content",
	})

	fs := &tests.CreateErrorFS{Filesystem: memfs.New()}
	if err := Tar(buf, fs); err == nil {
		t.Fatal("expected error when file create fails")
	}
}

func TestTar_CopyError(t *testing.T) {
	buf := buildTar(t, map[string]string{
		"values.yaml": "content",
	})

	fs := &tests.CopyErrorFS{Filesystem: memfs.New()}
	if err := Tar(buf, fs); err == nil {
		t.Fatal("expected error when file write fails")
	}
}

func TestTar_CloseError(t *testing.T) {
	buf := buildTar(t, map[string]string{
		"values.yaml": "content",
	})

	fs := &tests.CloseErrorFS{Filesystem: memfs.New()}
	if err := Tar(buf, fs); err == nil {
		t.Fatal("expected error when file close fails")
	}
}

func TestTar_ReadError(t *testing.T) {
	fs := memfs.New()
	if err := Tar(tests.ErrReader{}, fs); err == nil {
		t.Fatal("expected error when reader fails")
	}
}
