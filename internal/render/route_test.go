package render_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/katastroma/keleustes"
	"github.com/katastroma/orpheus/internal/render"
)

func TestRouter_Render(t *testing.T) {
	var r render.Router
	r.Register(keleustes.RendererType_RENDERER_TYPE_PLAIN, func(_ io.Reader) ([]byte, error) {
		return []byte("manifest"), nil
	})

	out, err := r.Render(keleustes.RendererType_RENDERER_TYPE_PLAIN, strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(out) != "manifest" {
		t.Errorf("expected %q, got %q", "manifest", string(out))
	}
}

func TestRouter_Render_NoMatch(t *testing.T) {
	var r render.Router
	r.Register(keleustes.RendererType_RENDERER_TYPE_PLAIN, func(_ io.Reader) ([]byte, error) {
		return nil, nil
	})

	if _, err := r.Render(keleustes.RendererType_RENDERER_TYPE_HELM, strings.NewReader("")); err == nil {
		t.Fatal("expected error when no backend registered for type")
	}
}

func TestRouter_Render_Empty(t *testing.T) {
	var r render.Router

	if _, err := r.Render(keleustes.RendererType_RENDERER_TYPE_PLAIN, strings.NewReader("")); err == nil {
		t.Fatal("expected error with no registered backends")
	}
}

func TestRouter_Render_BackendError(t *testing.T) {
	var r render.Router
	r.Register(keleustes.RendererType_RENDERER_TYPE_HELM, func(_ io.Reader) ([]byte, error) {
		return nil, fmt.Errorf("render failed")
	})

	if _, err := r.Render(keleustes.RendererType_RENDERER_TYPE_HELM, strings.NewReader("")); err == nil {
		t.Fatal("expected error from backend")
	}
}
