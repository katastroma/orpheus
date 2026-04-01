package serve_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	pb "github.com/katastroma/keleustes"
	"google.golang.org/grpc/metadata"

	"github.com/katastroma/orpheus/internal/serve"
	"github.com/katastroma/orpheus/internal/tests"
)

func renderContext(t *testing.T, rendererType pb.RendererType) context.Context {
	t.Helper()
	md := metadata.Pairs("renderer-type", rendererType.String())
	return metadata.NewIncomingContext(t.Context(), md)
}

func TestRender(t *testing.T) {
	var forwarded []byte

	svc := serve.New(
		slog.Default(),
		func(_ pb.RendererType, _ io.Reader) ([]byte, error) {
			return []byte("kind: Deployment"), nil
		},
		func(_ context.Context, manifest []byte) error {
			forwarded = manifest
			return nil
		},
	)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(forwarded) != "kind: Deployment" {
		t.Errorf("expected %q, got %q", "kind: Deployment", string(forwarded))
	}

	if len(stream.Responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(stream.Responses))
	}
}

func TestRender_MissingMetadata(t *testing.T) {
	svc := serve.New(slog.Default(), nil, nil)

	stream := &tests.MockRenderServer{
		Ctx: t.Context(),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when metadata is missing")
	}
}

func TestRender_RenderError(t *testing.T) {
	svc := serve.New(
		slog.Default(),
		func(_ pb.RendererType, _ io.Reader) ([]byte, error) {
			return nil, fmt.Errorf("render failed")
		},
		nil,
	)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when render fails")
	}
}

func TestRender_ForwardError(t *testing.T) {
	svc := serve.New(
		slog.Default(),
		func(_ pb.RendererType, _ io.Reader) ([]byte, error) {
			return []byte("manifest"), nil
		},
		func(_ context.Context, _ []byte) error {
			return fmt.Errorf("forward failed")
		},
	)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when forward fails")
	}
}

func TestRender_SendError(t *testing.T) {
	svc := serve.New(
		slog.Default(),
		func(_ pb.RendererType, _ io.Reader) ([]byte, error) {
			return []byte("manifest"), nil
		},
		func(_ context.Context, _ []byte) error { return nil },
	)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		SendErr:  fmt.Errorf("send failed"),
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when send fails")
	}
}
