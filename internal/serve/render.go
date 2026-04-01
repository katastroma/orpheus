//revive:disable:package-comments
package serve

import (
	"context"
	"fmt"

	pb "github.com/katastroma/keleustes"
	"google.golang.org/grpc/metadata"
)

// rendererTypeKey is the gRPC metadata key carrying the renderer type.
const rendererTypeKey = "renderer-type"

// Render receives a tar archive from the stream, dispatches to the
// appropriate rendering backend, and streams the result to the orderer.
func (s *Service) Render(stream pb.RendererService_RenderServer) error {
	ctx := stream.Context()
	s.log.InfoContext(ctx, "render requested")

	rendererType, err := readRendererType(ctx)
	if err != nil {
		s.log.ErrorContext(ctx, "reading renderer type failed", "error", err)
		return fmt.Errorf("reading renderer type: %w", err)
	}
	s.log.InfoContext(ctx, "renderer type received", "renderer", rendererType.String())

	s.log.DebugContext(ctx, "rendering manifests")
	manifest, err := s.renderFn(rendererType, &streamReader{stream: stream})
	if err != nil {
		s.log.ErrorContext(ctx, "rendering failed", "error", err)
		return fmt.Errorf("rendering: %w", err)
	}
	s.log.DebugContext(ctx, "manifests rendered", "bytes", len(manifest))

	s.log.DebugContext(ctx, "streaming to orderer")
	if err = s.streamFn(ctx, manifest); err != nil {
		s.log.ErrorContext(ctx, "streaming to orderer failed", "error", err)
		return fmt.Errorf("streaming to orderer: %w", err)
	}
	s.log.InfoContext(ctx, "streamed to orderer")

	if err = stream.SendAndClose(&pb.RenderResponse{}); err != nil {
		s.log.ErrorContext(ctx, "sending response failed", "error", err)
		return fmt.Errorf("sending response: %w", err)
	}

	return nil
}

func readRendererType(ctx context.Context) (pb.RendererType, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return pb.RendererType_RENDERER_TYPE_UNSPECIFIED, fmt.Errorf("no gRPC metadata")
	}

	values := md.Get(rendererTypeKey)
	if len(values) == 0 {
		return pb.RendererType_RENDERER_TYPE_UNSPECIFIED, fmt.Errorf("missing %s metadata", rendererTypeKey)
	}

	rendererType, ok := pb.RendererType_value[values[0]]
	if !ok {
		return pb.RendererType_RENDERER_TYPE_UNSPECIFIED, fmt.Errorf("unknown renderer type %q", values[0])
	}

	return pb.RendererType(rendererType), nil
}
