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
	backend := &tests.MockBackend{RenderResult: []byte("kind: Deployment")}

	svc := serve.New(
		slog.Default(),
		routerWith(backend),
		func(_ context.Context, _ []byte) error { return nil },
		32*1024,
	)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stream.Responses) != 1 {
		t.Fatalf("expected 1 response, got %d", len(stream.Responses))
	}
}

func TestRender_SendAndCloseError(t *testing.T) {
	backend := &tests.MockBackend{RenderResult: []byte("kind: Deployment")}

	svc := serve.New(
		slog.Default(),
		routerWith(backend),
		func(_ context.Context, _ []byte) error { return nil },
		32*1024,
	)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		SendErr:  fmt.Errorf("send failed"),
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when SendAndClose fails")
	}
}

func TestRender_MissingMetadata(t *testing.T) {
	svc := serve.New(slog.Default(), &render.Router{}, nil, 32*1024)

	stream := &tests.MockRenderServer{
		Ctx: t.Context(),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when metadata is missing")
	}
}

func TestRender_BackendLookupError(t *testing.T) {
	svc := serve.New(slog.Default(), &render.Router{}, nil, 32*1024)

	stream := &tests.MockRenderServer{
		Ctx: renderContext(t, pb.RendererType_RENDERER_TYPE_HELM),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when backend not registered")
	}
}

func TestRender_ReceiveError(t *testing.T) {
	backend := &tests.MockBackend{ReceiveErr: fmt.Errorf("receive failed")}

	svc := serve.New(slog.Default(), routerWith(backend), nil, 32*1024)

	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when receive fails")
	}
}

func TestRenderStream(t *testing.T) {
	backend := &tests.MockBackend{RenderResult: []byte("kind: Deployment")}

	svc := serve.New(slog.Default(), routerWith(backend), nil, 32*1024)

	stream := &tests.MockRenderStreamServer{
		Requests: []*pb.RenderStreamRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.RenderStream(stream); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stream.Sent) == 0 {
		t.Fatal("expected at least one response chunk")
	}

	var got []byte
	for _, resp := range stream.Sent {
		got = append(got, resp.GetData()...)
	}
	if string(got) != "kind: Deployment" {
		t.Errorf("expected %q, got %q", "kind: Deployment", string(got))
	}
}

func TestRenderStream_MissingMetadata(t *testing.T) {
	svc := serve.New(slog.Default(), &render.Router{}, nil, 32*1024)

	stream := &tests.MockRenderStreamServer{
		Ctx: t.Context(),
	}

	if err := svc.RenderStream(stream); err == nil {
		t.Fatal("expected error when metadata is missing")
	}
}

func TestRenderStream_BackendLookupError(t *testing.T) {
	svc := serve.New(slog.Default(), &render.Router{}, nil, 32*1024)

	stream := &tests.MockRenderStreamServer{
		Ctx: renderContext(t, pb.RendererType_RENDERER_TYPE_HELM),
	}

	if err := svc.RenderStream(stream); err == nil {
		t.Fatal("expected error when backend not registered")
	}
}

func TestRenderStream_ReceiveError(t *testing.T) {
	backend := &tests.MockBackend{ReceiveErr: fmt.Errorf("receive failed")}

	svc := serve.New(slog.Default(), routerWith(backend), nil, 32*1024)

	stream := &tests.MockRenderStreamServer{
		Requests: []*pb.RenderStreamRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.RenderStream(stream); err == nil {
		t.Fatal("expected error when receive fails")
	}
}

func TestRenderStream_RenderError(t *testing.T) {
	backend := &tests.MockBackend{RenderErr: fmt.Errorf("render failed")}

	svc := serve.New(slog.Default(), routerWith(backend), nil, 32*1024)

	stream := &tests.MockRenderStreamServer{
		Requests: []*pb.RenderStreamRequest{{Data: []byte("tar data")}},
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.RenderStream(stream); err == nil {
		t.Fatal("expected error when render fails")
	}
}

func TestRenderStream_SendError(t *testing.T) {
	backend := &tests.MockBackend{RenderResult: []byte("kind: Deployment")}

	svc := serve.New(slog.Default(), routerWith(backend), nil, 32*1024)

	stream := &tests.MockRenderStreamServer{
		Requests: []*pb.RenderStreamRequest{{Data: []byte("tar data")}},
		SendErr:  fmt.Errorf("send failed"),
		Ctx:      renderContext(t, pb.RendererType_RENDERER_TYPE_PLAIN),
	}

	if err := svc.RenderStream(stream); err == nil {
		t.Fatal("expected error when send fails")
	}
}
