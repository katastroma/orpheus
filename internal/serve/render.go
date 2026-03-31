//revive:disable:package-comments
package serve

import (
	"fmt"

	pb "github.com/katastroma/keleustes"

	"github.com/katastroma/orpheus/internal/extract"
)

// Render receives a tar archive from the stream, extracts it, renders
// manifests, and streams them to the orderer.
func (s *Service) Render(stream pb.RendererService_RenderServer) error {
	ctx := stream.Context()
	s.log.InfoContext(ctx, "render requested")

	s.log.DebugContext(ctx, "extracting source archive")
	files, err := extract.Tar(&streamReader{stream: stream})
	if err != nil {
		s.log.ErrorContext(ctx, "extraction failed", "error", err)
		return fmt.Errorf("extracting source: %w", err)
	}
	s.log.DebugContext(ctx, "source extracted", "files", len(files))

	s.log.DebugContext(ctx, "rendering manifests")
	manifest, err := s.renderFn(files)
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
