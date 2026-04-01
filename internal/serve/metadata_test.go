package serve

import (
	"testing"

	"github.com/katastroma/keleustes"
	"google.golang.org/grpc/metadata"
)

func TestReadRendererType(t *testing.T) {
	md := metadata.Pairs(keleustes.RendererTypeMetadataKey, keleustes.RendererType_RENDERER_TYPE_HELM.String())
	ctx := metadata.NewIncomingContext(t.Context(), md)

	got, err := readRendererType(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != keleustes.RendererType_RENDERER_TYPE_HELM {
		t.Errorf("expected HELM, got %v", got)
	}
}

func TestReadRendererType_NoMetadata(t *testing.T) {
	if _, err := readRendererType(t.Context()); err == nil {
		t.Fatal("expected error when no metadata present")
	}
}

func TestReadRendererType_MissingKey(t *testing.T) {
	md := metadata.Pairs("other-key", "value")
	ctx := metadata.NewIncomingContext(t.Context(), md)

	if _, err := readRendererType(ctx); err == nil {
		t.Fatal("expected error when renderer-type key is missing")
	}
}

func TestReadRendererType_UnknownType(t *testing.T) {
	md := metadata.Pairs(keleustes.RendererTypeMetadataKey, "RENDERER_TYPE_UNKNOWN")
	ctx := metadata.NewIncomingContext(t.Context(), md)

	if _, err := readRendererType(ctx); err == nil {
		t.Fatal("expected error for unknown renderer type")
	}
}
