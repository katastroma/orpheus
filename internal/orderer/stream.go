//revive:disable:package-comments
package orderer

import (
	"context"
	"fmt"

	pb "github.com/katastroma/diataxis"
	"google.golang.org/grpc"
)

// StreamFunc streams a rendered manifest blob to the orderer.
type StreamFunc func(ctx context.Context, manifest []byte) error

// NewStreamFunc returns a StreamFunc that opens an Order stream on conn
// and sends the manifest.
func NewStreamFunc(conn grpc.ClientConnInterface) StreamFunc {
	return func(ctx context.Context, manifest []byte) error {
		stream, err := pb.NewOrdererServiceClient(conn).Order(ctx)
		if err != nil {
			return fmt.Errorf("opening order stream: %w", err)
		}

		if err = stream.Send(&pb.OrderRequest{Manifest: manifest}); err != nil {
			return fmt.Errorf("sending manifest to orderer: %w", err)
		}

		return stream.CloseSend()
	}
}
