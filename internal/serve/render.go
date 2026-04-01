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

	rendererType, err := readRendererType(ctx)
	if err != nil {
		fail(ctx, s.log, "reading renderer type failed", err)
		return fmt.Errorf("reading renderer type: %w", err)
	}

	log := s.log.With("renderer", rendererType.String())
	log.InfoContext(ctx, "renderer type received")

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

	log.DebugContext(ctx, "rendering manifests")
	manifest, err := backend.Render()
	if err != nil {
		fail(ctx, log, "rendering failed", err)
		return nil
	}
	log.DebugContext(ctx, "manifests rendered", "bytes", len(manifest))

	log.DebugContext(ctx, "streaming to orderer")
	if err = s.streamFn(ctx, manifest); err != nil {
		fail(ctx, log, "streaming to orderer failed", err)
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
