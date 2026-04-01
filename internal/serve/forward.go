//revive:disable:package-comments
package serve

import (
	"context"
	"log/slog"

	"github.com/katastroma/orpheus/internal/orderer"
	"github.com/katastroma/orpheus/internal/render"
)

func renderAndForward(ctx context.Context, log *slog.Logger, backend render.Backend, streamFn orderer.StreamFunc) {
	ctx = context.WithoutCancel(ctx)

	log.DebugContext(ctx, "rendering manifests")
	manifest, err := backend.Render()
	if err != nil {
		fail(ctx, log, "rendering failed", err)
		return
	}
	log.DebugContext(ctx, "manifests rendered", "bytes", len(manifest))

	log.DebugContext(ctx, "streaming to orderer")
	if err = streamFn(ctx, manifest); err != nil {
		fail(ctx, log, "streaming to orderer failed", err)
		return
	}
	log.InfoContext(ctx, "streamed to orderer")
}
