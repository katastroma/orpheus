//revive:disable:package-comments
package orderer

import (
	"context"
	"fmt"

	pb "github.com/katastroma/diataxis"
	"google.golang.org/grpc"
)

// StreamFunc streams rendered manifests to the orderer.
type StreamFunc func(ctx context.Context, manifests [][]byte) error

// NewStreamFunc returns a StreamFunc that opens an Order stream on conn
// and sends each manifest.
func NewStreamFunc(conn grpc.ClientConnInterface) StreamFunc {
	return func(ctx context.Context, manifests [][]byte) error {
		stream, err := pb.NewOrdererServiceClient(conn).Order(ctx)
		if err != nil {
			return fmt.Errorf("opening order stream: %w", err)
		}

		for _, m := range manifests {
			if err = stream.Send(&pb.OrderRequest{Manifest: m}); err != nil {
				return fmt.Errorf("sending manifest to orderer: %w", err)
			}
		}

		return stream.CloseSend()
	}
}
