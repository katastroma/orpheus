package plain

import (
	"strings"
	"testing"

	"github.com/katastroma/orpheus/internal/tests"
)

func TestReceive_TarReadError(t *testing.T) {
	b := &Backend{}
	if err := b.Receive(tests.ErrReader{}); err == nil {
		t.Fatal("expected error when tar read fails")
	}
}

func TestReceive_TruncatedEntry(t *testing.T) {
	b := &Backend{}
	if err := b.Receive(tests.TruncatedTarReader(t, "app.yaml", 1000)); err == nil {
		t.Fatal("expected error for truncated tar entry")
	}
}

func TestReceive_SkipsDirectories(t *testing.T) {
	b := &Backend{}
	r := tests.TarReaderWithDir(t, "deploy/", map[string][]byte{
		"deploy/app.yaml": []byte("kind: Pod"),
	})

	if err := b.Receive(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, err := b.Render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}

	if !strings.Contains(string(out), "kind: Pod") {
		t.Errorf("expected Pod in output, got %q", string(out))
	}
}

func TestReceive_CorruptArchive(t *testing.T) {
	b := &Backend{}
	if err := b.Receive(strings.NewReader("not a tar")); err == nil {
		t.Fatal("expected error for corrupt archive")
	}
}
