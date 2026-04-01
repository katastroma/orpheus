package render_test

import (
	"testing"

	"github.com/katastroma/keleustes"
	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/tests"
)

func TestRouter_Lookup(t *testing.T) {
	var r render.Router
	r.Register(keleustes.RendererType_RENDERER_TYPE_PLAIN, &tests.MockBackend{})

	backend, err := r.Lookup(keleustes.RendererType_RENDERER_TYPE_PLAIN)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if backend == nil {
		t.Fatal("expected non-nil backend")
	}
}

func TestRouter_Lookup_NotRegistered(t *testing.T) {
	var r render.Router
	r.Register(keleustes.RendererType_RENDERER_TYPE_PLAIN, &tests.MockBackend{})

	if _, err := r.Lookup(keleustes.RendererType_RENDERER_TYPE_HELM); err == nil {
		t.Fatal("expected error when no backend registered for type")
	}
}

func TestRouter_Lookup_Empty(t *testing.T) {
	var r render.Router

	if _, err := r.Lookup(keleustes.RendererType_RENDERER_TYPE_PLAIN); err == nil {
		t.Fatal("expected error with no registered backends")
	}
}
