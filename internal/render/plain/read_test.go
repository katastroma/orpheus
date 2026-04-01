package plain

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/tests"
)

func TestRender_TarReadError(t *testing.T) {
	if _, err := Render(tests.ErrReader{}); err == nil {
		t.Fatal("expected error when tar read fails")
	}
}

func TestRender_TruncatedEntry(t *testing.T) {
	r := tests.TruncatedTarReader(t, "app.yaml", 1000)

	if _, err := Render(r); err == nil {
		t.Fatal("expected error for truncated tar entry")
	}
}

func TestRender_SkipsDirectories(t *testing.T) {
	r := tests.TarReaderWithDir(t, "deploy/", map[string][]byte{
		"deploy/app.yaml": []byte("kind: Pod"),
	})

	out, err := Render(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(string(out), "kind: Pod") {
		t.Errorf("expected Pod in output, got %q", string(out))
	}
}

func TestRender_CorruptArchive(t *testing.T) {
	if _, err := Render(strings.NewReader("not a tar")); err == nil {
		t.Fatal("expected error for corrupt archive")
	}
}
