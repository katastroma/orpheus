package orderer

import (
	"fmt"
	"log/slog"
	"testing"

	"google.golang.org/grpc"

	"github.com/katastroma/orpheus/internal/tests"
)

func TestSend_OpenError(t *testing.T) {
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) {
			return nil, fmt.Errorf("connection refused")
		},
	}

	if err := send(t.Context(), slog.Default(), nil, conn); err == nil {
		t.Fatal("expected error when stream fails to open")
	}
}

func TestSend_SendError(t *testing.T) {
	cs := &tests.MockClientStream{SendMsgErr: fmt.Errorf("send failed"), Ctx: t.Context()}
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) { return cs, nil },
	}

	if err := send(t.Context(), slog.Default(), nil, conn); err == nil {
		t.Fatal("expected error when send fails")
	}
}

func TestSend_CloseAndRecvError(t *testing.T) {
	cs := &tests.MockClientStream{RecvMsgErr: fmt.Errorf("server error"), Ctx: t.Context()}
	conn := &tests.MockClientConn{
		NewStreamFn: func() (grpc.ClientStream, error) { return cs, nil },
	}

	if err := send(t.Context(), slog.Default(), nil, conn); err == nil {
		t.Fatal("expected error when CloseAndRecv fails")
	}
}
