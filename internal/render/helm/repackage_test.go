package helm

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"testing"

	"github.com/katastroma/orpheus/internal/tests"
)

func TestRepackage(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"Chart.yaml": []byte("name: test"),
	})

	result := repackage(r)

	gr, err := gzip.NewReader(result)
	if err != nil {
		t.Fatalf("expected valid gzip: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	hdr, err := tr.Next()
	if err != nil {
		t.Fatalf("expected tar entry: %v", err)
	}

	if hdr.Name != directoryPrefix+"Chart.yaml" {
		t.Errorf("expected %q, got %q", directoryPrefix+"Chart.yaml", hdr.Name)
	}
}

func TestRewriteEntries_WriteHeaderError(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"Chart.yaml": []byte("name: test"),
	})

	tr := tar.NewReader(r)
	tw := tar.NewWriter(tests.ErrWriter{})

	if err := rewriteEntries(tr, tw); err == nil {
		t.Fatal("expected error when write header fails")
	}
}

func TestRewriteEntries_CopyError(t *testing.T) {
	r := tests.TarReader(t, map[string][]byte{
		"Chart.yaml": []byte("name: test"),
	})

	tr := tar.NewReader(r)
	tw := tar.NewWriter(&tests.FailAfterNWriter{N: 512})

	if err := rewriteEntries(tr, tw); err == nil {
		t.Fatal("expected error when copy fails")
	}
}

func TestRepackage_ReadError(t *testing.T) {
	result := repackage(tests.ErrReader{})

	if _, err := io.ReadAll(result); err == nil {
		t.Fatal("expected error when source reader fails")
	}
}
