//revive:disable:package-comments
package orderer

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

// StreamFunc streams a rendered manifest blob to the orderer.
type StreamFunc func(ctx context.Context, manifest []byte) error

// NewStreamFunc returns a StreamFunc that opens an Order stream on conn
// and sends the manifest.
func NewStreamFunc(log *slog.Logger, conn grpc.ClientConnInterface) StreamFunc {
	return func(ctx context.Context, manifest []byte) error {
		return send(ctx, log, manifest, conn)
	}
}
