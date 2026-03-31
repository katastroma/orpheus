package serve

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"

	pb "github.com/katastroma/keleustes"

	"github.com/katastroma/orpheus/internal/render"
	"github.com/katastroma/orpheus/internal/tests"
)

func tarFromFiles(t *testing.T, files render.Files) []byte {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for name, content := range files {
		if err := tw.WriteHeader(&tar.Header{
			Name:     name,
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatalf("writing tar header for %s: %v", name, err)
		}

		if _, err := tw.Write(content); err != nil {
			t.Fatalf("writing tar content for %s: %v", name, err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar: %v", err)
	}

	return buf.Bytes()
}

func successRender(_ render.Files) ([]byte, error) {
	return []byte("kind: Deployment"), nil
}

func errorRender(_ render.Files) ([]byte, error) {
	return nil, fmt.Errorf("render failed")
}

func successForward(_ context.Context, _ []byte) error {
	return nil
}

func errorForward(_ context.Context, _ []byte) error {
	return fmt.Errorf("forward failed")
}

func TestRender(t *testing.T) {
	data := tarFromFiles(t, render.Files{"app.yaml": []byte("kind: Deployment")})

	var forwarded []byte
	captureFn := func(_ context.Context, manifest []byte) error {
		forwarded = manifest
		return nil
	}

	svc := New(slog.Default(), successRender, captureFn)
	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: data}},
		Ctx:      t.Context(),
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

func TestRender_ExtractError(t *testing.T) {
	svc := New(slog.Default(), successRender, successForward)
	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("not a tar")}},
		Ctx:      t.Context(),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error for corrupt tar")
	}
}

func TestRender_RenderError(t *testing.T) {
	data := tarFromFiles(t, render.Files{"app.yaml": []byte("kind: Pod")})

	svc := New(slog.Default(), errorRender, successForward)
	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: data}},
		Ctx:      t.Context(),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when render fails")
	}
}

func TestRender_ForwardError(t *testing.T) {
	data := tarFromFiles(t, render.Files{"app.yaml": []byte("kind: Pod")})

	svc := New(slog.Default(), successRender, errorForward)
	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: data}},
		Ctx:      t.Context(),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when forward fails")
	}
}

func TestRender_ReadError(t *testing.T) {
	svc := New(slog.Default(), successRender, successForward)
	stream := &tests.MockRenderServer{
		RecvErr: fmt.Errorf("recv failed"),
		Ctx:     t.Context(),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when read fails")
	}
}

func TestRender_SendError(t *testing.T) {
	data := tarFromFiles(t, render.Files{"app.yaml": []byte("kind: Pod")})

	svc := New(slog.Default(), successRender, successForward)
	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: data}},
		SendErr:  fmt.Errorf("send failed"),
		Ctx:      t.Context(),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when send fails")
	}
}
