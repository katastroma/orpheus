package serve

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
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

	svc := New(successRender, captureFn)
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
}

func TestRender_ExtractError(t *testing.T) {
	svc := New(successRender, successForward)
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

	svc := New(errorRender, successForward)
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

	svc := New(successRender, errorForward)
	stream := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: data}},
		Ctx:      t.Context(),
	}

	if err := svc.Render(stream); err == nil {
		t.Fatal("expected error when forward fails")
	}
}
