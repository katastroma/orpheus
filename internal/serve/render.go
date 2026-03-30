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
	files, err := extract.Tar(&streamReader{stream: stream})
	if err != nil {
		return fmt.Errorf("extracting source: %w", err)
	}

	manifests, err := s.renderFn(files)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	if err = s.streamFn(stream.Context(), manifests); err != nil {
		return fmt.Errorf("streaming to orderer: %w", err)
	}

	return nil
}
