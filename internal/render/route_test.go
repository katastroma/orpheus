package render_test

import (
	"fmt"
	"testing"

	"github.com/katastroma/orpheus/internal/render"
)

func TestRouter_Render(t *testing.T) {
	var r render.Router
	r.Register("test", func(_ render.Files) bool { return true }, func(_ render.Files) ([]byte, error) {
		return []byte("manifest"), nil
	})

	out, err := r.Render(render.Files{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(out) != "manifest" {
		t.Errorf("expected %q, got %q", "manifest", string(out))
	}
}

func TestRouter_Render_Priority(t *testing.T) {
	var r render.Router
	r.Register("first", func(_ render.Files) bool { return true }, func(_ render.Files) ([]byte, error) {
		return []byte("first"), nil
	})
	r.Register("second", func(_ render.Files) bool { return true }, func(_ render.Files) ([]byte, error) {
		return []byte("second"), nil
	})

	out, err := r.Render(render.Files{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(out) != "first" {
		t.Errorf("expected first registered backend to win, got %q", string(out))
	}
}

func TestRouter_Render_SkipsNonMatching(t *testing.T) {
	var r render.Router
	r.Register("no", func(_ render.Files) bool { return false }, func(_ render.Files) ([]byte, error) {
		return []byte("wrong"), nil
	})
	r.Register("yes", func(_ render.Files) bool { return true }, func(_ render.Files) ([]byte, error) {
		return []byte("right"), nil
	})

	out, err := r.Render(render.Files{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(out) != "right" {
		t.Errorf("expected matching backend, got %q", string(out))
	}
}

func TestRouter_Render_NoMatch(t *testing.T) {
	var r render.Router
	r.Register("no", func(_ render.Files) bool { return false }, func(_ render.Files) ([]byte, error) {
		return nil, nil
	})

	if _, err := r.Render(render.Files{}); err == nil {
		t.Fatal("expected error when no backend matches")
	}
}

func TestRouter_Render_Empty(t *testing.T) {
	var r render.Router

	if _, err := r.Render(render.Files{}); err == nil {
		t.Fatal("expected error with no registered backends")
	}
}

func TestRouter_Render_BackendError(t *testing.T) {
	var r render.Router
	r.Register("fail", func(_ render.Files) bool { return true }, func(_ render.Files) ([]byte, error) {
		return nil, fmt.Errorf("render failed")
	})

	if _, err := r.Render(render.Files{}); err == nil {
		t.Fatal("expected error from backend")
	}
}
