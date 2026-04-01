package serve

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/katastroma/orpheus/internal/tests"
)

func TestRenderAndForward(t *testing.T) {
	var forwarded []byte
	backend := &tests.MockBackend{RenderResult: []byte("kind: Deployment")}

	renderAndForward(t.Context(), slog.Default(), backend, func(_ context.Context, manifest []byte) error {
		forwarded = manifest
		return nil
	})

	if string(forwarded) != "kind: Deployment" {
		t.Errorf("expected %q, got %q", "kind: Deployment", string(forwarded))
	}
}

func TestRenderAndForward_RenderError(t *testing.T) {
	backend := &tests.MockBackend{RenderErr: fmt.Errorf("render failed")}

	renderAndForward(t.Context(), slog.Default(), backend, func(_ context.Context, _ []byte) error {
		t.Fatal("streamFn should not be called when render fails")
		return nil
	})
}

func TestRenderAndForward_ForwardError(t *testing.T) {
	backend := &tests.MockBackend{RenderResult: []byte("manifest")}

	renderAndForward(t.Context(), slog.Default(), backend, func(_ context.Context, _ []byte) error {
		return fmt.Errorf("forward failed")
	})
}
