package serve_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	pb "github.com/katastroma/keleustes"
	"google.golang.org/grpc/metadata"

	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/serve"
	"github.com/katastroma/orpheus/internal/tests"
)

func renderContext(t *testing.T, rendererType pb.RendererType) context.Context {
	t.Helper()
	md := metadata.Pairs(pb.RendererTypeMetadataKey, rendererType.String())
	return metadata.NewIncomingContext(t.Context(), md)
}

func routerWith(backend *tests.MockBackend) *render.Router {
	var r render.Router
	r.Register(pb.RendererType_RENDERER_TYPE_PLAIN, backend)
	return &r
}

func TestRender(t *testing.T) {
	var forwarded []byte
	backend := &tests.MockBackend{RenderResult: []byte("kind: Deployment")}

	svc := serve.New(
		slog.Default(),
		routerWith(backend),
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
	svc := serve.New(slog.Default(), &render.Router{}, nil)

	stream := &tests.MockRenderServer{
		Ctx: t.Context(),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when metadata is missing")
	}
}

func TestRender_BackendLookupError(t *testing.T) {
	svc := serve.New(slog.Default(), &render.Router{}, nil)

	stream := &tests.MockRenderServer{
		Ctx: renderContext(t, pb.RendererType_RENDERER_TYPE_HELM),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when backend not registered")
	}
}

func TestRender_ReceiveError(t *testing.T) {
	backend := &tests.MockBackend{ReceiveErr: fmt.Errorf("receive failed")}

	svc := serve.New(slog.Default(), routerWith(backend), nil)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when receive fails")
	}
}

func TestRender_RenderError(t *testing.T) {
	backend := &tests.MockBackend{RenderErr: fmt.Errorf("render failed")}

	svc := serve.New(slog.Default(), routerWith(backend), nil)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err != nil {
		t.Fatal("expected nil return after SendAndClose even when render fails")
	}

	if len(stream.Responses) != 1 {
		t.Fatalf("expected response sent before render failure, got %d", len(stream.Responses))
	}
}

func TestRender_ForwardError(t *testing.T) {
	backend := &tests.MockBackend{RenderResult: []byte("manifest")}

	svc := serve.New(
		slog.Default(),
		routerWith(backend),
		func(_ context.Context, _ []byte) error {
			return fmt.Errorf("forward failed")
		},
	)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err != nil {
		t.Fatal("expected nil return after SendAndClose even when forward fails")
	}
}

func TestRender_SendError(t *testing.T) {
	backend := &tests.MockBackend{RenderResult: []byte("manifest")}

	svc := serve.New(
		slog.Default(),
		routerWith(backend),
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
