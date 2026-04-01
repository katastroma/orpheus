//revive:disable:package-comments
package tests

import (
	"archive/tar"
	"bytes"
	"testing"
)

// TarReader builds a tar archive from a file map and returns it as a reader.
func TarReader(t *testing.T, files map[string][]byte) *bytes.Reader {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for name, data := range files {
		if err := tw.WriteHeader(&tar.Header{
			Name:     name,
			Size:     int64(len(data)),
			Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatalf("writing tar header for %s: %v", name, err)
		}

		if _, err := tw.Write(data); err != nil {
			t.Fatalf("writing tar content for %s: %v", name, err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar: %v", err)
	}

	return bytes.NewReader(buf.Bytes())
}

// TarReaderWithDir builds a tar archive containing a directory entry followed
// by file entries from the map.
func TarReaderWithDir(t *testing.T, dir string, files map[string][]byte) *bytes.Reader {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	if err := tw.WriteHeader(&tar.Header{
		Name:     dir,
		Typeflag: tar.TypeDir,
	}); err != nil {
		t.Fatalf("writing dir header: %v", err)
	}

	for name, data := range files {
		if err := tw.WriteHeader(&tar.Header{
			Name:     name,
			Size:     int64(len(data)),
			Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatalf("writing tar header for %s: %v", name, err)
		}

		if _, err := tw.Write(data); err != nil {
			t.Fatalf("writing tar content for %s: %v", name, err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar: %v", err)
	}

	return bytes.NewReader(buf.Bytes())
}

// TruncatedTarReader builds a tar with a header claiming size bytes but no
// content, producing a truncated entry that causes read errors.
func TruncatedTarReader(t *testing.T, name string, size int64) *bytes.Reader {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	if err := tw.WriteHeader(&tar.Header{
		Name:     name,
		Size:     size,
		Typeflag: tar.TypeReg,
	}); err != nil {
		t.Fatalf("writing tar header: %v", err)
	}

	return bytes.NewReader(buf.Bytes())
}
