//revive:disable:package-comments
package serve

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/katastroma/keleustes"
	"google.golang.org/grpc/metadata"
)

func fail(ctx context.Context, log *slog.Logger, msg string, err error) {
	log.ErrorContext(ctx, msg, "error", err)
}

// Render receives a tar archive from the stream, dispatches to the
// appropriate rendering backend, and streams the result to the orderer.
func (s *Service) Render(stream pb.RendererService_RenderServer) error {
	ctx := stream.Context()
	s.log.InfoContext(ctx, "render requested")

	s.log.DebugContext(ctx, "looking up renderer")
	rendererType, err := readRendererType(ctx)
	if err != nil {
		fail(ctx, s.log, "reading renderer type failed", err)
		return fmt.Errorf("reading renderer type: %w", err)
	}
	log := s.log.With("renderer", rendererType.String())
	log.DebugContext(ctx, "renderer type received")

	log.DebugContext(ctx, "looking up backend")
	backend, err := s.router.Lookup(rendererType)
	if err != nil {
		fail(ctx, log, "backend lookup failed", err)
		return fmt.Errorf("backend lookup: %w", err)
	}
	log.DebugContext(ctx, "backend found")

	log.DebugContext(ctx, "receiving source content")
	if err = backend.Receive(&streamReader{stream: stream}); err != nil {
		fail(ctx, log, "receiving failed", err)
		return fmt.Errorf("receiving: %w", err)
	}
	log.DebugContext(ctx, "source content received")

	if err = stream.SendAndClose(&pb.RenderResponse{}); err != nil {
		fail(ctx, log, "sending response failed", err)
		return fmt.Errorf("sending response: %w", err)
	}

	go renderAndForward(ctx, log, backend, s.streamFn)

	return nil
}

// RenderStream receives a tar archive over the bidi stream, renders it,
// and streams the resulting manifest back in chunks.
func (s *Service) RenderStream(stream pb.RendererService_RenderStreamServer) error {
	ctx := stream.Context()
	s.log.InfoContext(ctx, "render stream requested")

	rendererType, err := readRendererType(ctx)
	if err != nil {
		fail(ctx, s.log, "reading renderer type failed", err)
		return fmt.Errorf("reading renderer type: %w", err)
	}
	log := s.log.With("renderer", rendererType.String())

	backend, err := s.router.Lookup(rendererType)
	if err != nil {
		fail(ctx, log, "backend lookup failed", err)
		return fmt.Errorf("backend lookup: %w", err)
	}

	if err = backend.Receive(&renderStreamReader{stream: stream}); err != nil {
		fail(ctx, log, "receiving failed", err)
		return fmt.Errorf("receiving: %w", err)
	}

	manifest, err := backend.Render()
	if err != nil {
		fail(ctx, log, "rendering failed", err)
		return fmt.Errorf("rendering: %w", err)
	}

	for len(manifest) > 0 {
		n := s.chunkSize
		if n > len(manifest) {
			n = len(manifest)
		}
		if err = stream.Send(&pb.RenderStreamResponse{Data: manifest[:n]}); err != nil {
			fail(ctx, log, "sending chunk failed", err)
			return fmt.Errorf("sending chunk: %w", err)
		}
		manifest = manifest[n:]
	}

	return nil
}

func readRendererType(ctx context.Context) (pb.RendererType, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return pb.RendererType_RENDERER_TYPE_UNSPECIFIED, fmt.Errorf("no gRPC metadata")
	}

	values := md.Get(pb.RendererTypeMetadataKey)
	if len(values) == 0 {
		return pb.RendererType_RENDERER_TYPE_UNSPECIFIED, fmt.Errorf("missing %s metadata", pb.RendererTypeMetadataKey)
	}

	rendererType, ok := pb.RendererType_value[values[0]]
	if !ok {
		return pb.RendererType_RENDERER_TYPE_UNSPECIFIED, fmt.Errorf("unknown renderer type %q", values[0])
	}

	return pb.RendererType(rendererType), nil
}
