package serve

import (
	"io"
	"testing"

	pb "github.com/katastroma/keleustes"

	"github.com/katastroma/orpheus/internal/tests"
)

func TestStreamReader_SingleMessage(t *testing.T) {
	mock := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("hello")}},
		Ctx:      t.Context(),
	}
	r := &streamReader{stream: mock}

	buf := make([]byte, 16)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(buf[:n]) != "hello" {
		t.Errorf("expected %q, got %q", "hello", string(buf[:n]))
	}
}

func TestStreamReader_MultipleMessages(t *testing.T) {
	mock := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{
			{Data: []byte("ab")},
			{Data: []byte("cd")},
		},
		Ctx: t.Context(),
	}
	r := &streamReader{stream: mock}

	result, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "abcd" {
		t.Errorf("expected %q, got %q", "abcd", string(result))
	}
}

func TestStreamReader_PartialRead(t *testing.T) {
	mock := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{{Data: []byte("abcdef")}},
		Ctx:      t.Context(),
	}
	r := &streamReader{stream: mock}

	buf := make([]byte, 3)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error on first read: %v", err)
	}
	if string(buf[:n]) != "abc" {
		t.Errorf("expected %q, got %q", "abc", string(buf[:n]))
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error on second read: %v", err)
	}
	if string(buf[:n]) != "def" {
		t.Errorf("expected %q, got %q", "def", string(buf[:n]))
	}
}

func TestStreamReader_EOF(t *testing.T) {
	mock := &tests.MockRenderServer{
		Requests: []*pb.RenderRequest{},
		Ctx:      t.Context(),
	}
	r := &streamReader{stream: mock}

	buf := make([]byte, 16)
	_, err := r.Read(buf)
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}

func TestRenderStreamReader_SingleMessage(t *testing.T) {
	mock := &tests.MockRenderStreamServer{
		Requests: []*pb.RenderStreamRequest{{Data: []byte("hello")}},
		Ctx:      t.Context(),
	}
	r := &renderStreamReader{stream: mock}

	buf := make([]byte, 16)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(buf[:n]) != "hello" {
		t.Errorf("expected %q, got %q", "hello", string(buf[:n]))
	}
}

func TestRenderStreamReader_MultipleMessages(t *testing.T) {
	mock := &tests.MockRenderStreamServer{
		Requests: []*pb.RenderStreamRequest{
			{Data: []byte("ab")},
			{Data: []byte("cd")},
		},
		Ctx: t.Context(),
	}
	r := &renderStreamReader{stream: mock}

	result, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "abcd" {
		t.Errorf("expected %q, got %q", "abcd", string(result))
	}
}

func TestRenderStreamReader_PartialRead(t *testing.T) {
	mock := &tests.MockRenderStreamServer{
		Requests: []*pb.RenderStreamRequest{{Data: []byte("abcdef")}},
		Ctx:      t.Context(),
	}
	r := &renderStreamReader{stream: mock}

	buf := make([]byte, 3)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error on first read: %v", err)
	}
	if string(buf[:n]) != "abc" {
		t.Errorf("expected %q, got %q", "abc", string(buf[:n]))
	}

	n, err = r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error on second read: %v", err)
	}
	if string(buf[:n]) != "def" {
		t.Errorf("expected %q, got %q", "def", string(buf[:n]))
	}
}

func TestRenderStreamReader_EOF(t *testing.T) {
	mock := &tests.MockRenderStreamServer{
		Requests: []*pb.RenderStreamRequest{},
		Ctx:      t.Context(),
	}
	r := &renderStreamReader{stream: mock}

	buf := make([]byte, 16)
	_, err := r.Read(buf)
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}
