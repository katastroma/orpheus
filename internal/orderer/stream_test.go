//revive:disable:package-comments
package orderer_test

import (
	"log/slog"
	"testing"

	"github.com/katastroma/orpheus/internal/orderer"
	"github.com/katastroma/orpheus/internal/tests"
	"google.golang.org/grpc"
)

func TestNewStreamFunc(t *testing.T) {
	cs := &tests.MockClientStream{Ctx: t.Context()}
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) { return cs, nil },
	}
	streamFn := orderer.NewStreamFunc(slog.Default(), conn)

	if err := streamFn(t.Context(), []byte("kind: Deployment")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cs.SendMsgCount == 0 {
		t.Fatal("expected at least one message sent")
	}
}
