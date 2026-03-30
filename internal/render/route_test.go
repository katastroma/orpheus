package render

import (
	"fmt"
	"testing"
)

func TestRouter_Render(t *testing.T) {
	var r Router
	r.Register("test", func(_ Files) bool { return true }, func(_ Files) ([][]byte, error) {
		return [][]byte{[]byte("manifest")}, nil
	})

	manifests, err := r.Render(Files{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(manifests))
	}

	if string(manifests[0]) != "manifest" {
		t.Errorf("expected %q, got %q", "manifest", string(manifests[0]))
	}
}

func TestRouter_Render_Priority(t *testing.T) {
	var r Router
	r.Register("first", func(_ Files) bool { return true }, func(_ Files) ([][]byte, error) {
		return [][]byte{[]byte("first")}, nil
	})
	r.Register("second", func(_ Files) bool { return true }, func(_ Files) ([][]byte, error) {
		return [][]byte{[]byte("second")}, nil
	})

	manifests, err := r.Render(Files{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(manifests[0]) != "first" {
		t.Errorf("expected first registered backend to win, got %q", string(manifests[0]))
	}
}

func TestRouter_Render_SkipsNonMatching(t *testing.T) {
	var r Router
	r.Register("no", func(_ Files) bool { return false }, func(_ Files) ([][]byte, error) {
		return [][]byte{[]byte("wrong")}, nil
	})
	r.Register("yes", func(_ Files) bool { return true }, func(_ Files) ([][]byte, error) {
		return [][]byte{[]byte("right")}, nil
	})

	manifests, err := r.Render(Files{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(manifests[0]) != "right" {
		t.Errorf("expected matching backend, got %q", string(manifests[0]))
	}
}

func TestRouter_Render_NoMatch(t *testing.T) {
	var r Router
	r.Register("no", func(_ Files) bool { return false }, func(_ Files) ([][]byte, error) {
		return nil, nil
	})

	if _, err := r.Render(Files{}); err == nil {
		t.Fatal("expected error when no backend matches")
	}
}

func TestRouter_Render_Empty(t *testing.T) {
	var r Router

	if _, err := r.Render(Files{}); err == nil {
		t.Fatal("expected error with no registered backends")
	}
}

func TestRouter_Render_BackendError(t *testing.T) {
	var r Router
	r.Register("fail", func(_ Files) bool { return true }, func(_ Files) ([][]byte, error) {
		return nil, fmt.Errorf("render failed")
	})

	if _, err := r.Render(Files{}); err == nil {
		t.Fatal("expected error from backend")
	}
}
