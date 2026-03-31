//revive:disable:package-comments
package orderer

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/katastroma/diataxis"
	"google.golang.org/grpc"
)

func send(ctx context.Context, log *slog.Logger, manifest []byte, conn grpc.ClientConnInterface) error {
	stream, err := pb.NewOrdererServiceClient(conn).Order(ctx)
	if err != nil {
		return fmt.Errorf("opening order stream: %w", err)
	}

	if err = stream.Send(&pb.OrderRequest{Manifest: manifest}); err != nil {
		return fmt.Errorf("sending manifest to orderer: %w", err)
	}

	if _, err = stream.CloseAndRecv(); err != nil {
		log.ErrorContext(ctx, "orderer failed", "error", err)
		return fmt.Errorf("orderer: %w", err)
	}

	return nil
}
