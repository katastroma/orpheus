//revive:disable:package-comments
package serve

import pb "github.com/katastroma/keleustes"

// streamReader adapts a keleustes Render stream as an io.Reader. Each
// Read call consumes from the current message buffer, receiving the next
// RenderRequest when the buffer is exhausted.
type streamReader struct {
	stream pb.RendererService_RenderServer
	buf    []byte
}

// Read fills p from the stream, receiving new messages as needed.
func (r *streamReader) Read(p []byte) (int, error) {
	if len(r.buf) == 0 {
		req, err := r.stream.Recv()
		if err != nil {
			return 0, err
		}
		r.buf = req.GetData()
	}

	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

// renderStreamReader adapts a keleustes RenderStream bidi stream as an io.Reader.
// Each Read call consumes from the current message buffer, receiving the next
// RenderStreamRequest when the buffer is exhausted.
type renderStreamReader struct {
	stream pb.RendererService_RenderStreamServer
	buf    []byte
}

// Read fills p from the stream, receiving new messages as needed.
func (r *renderStreamReader) Read(p []byte) (int, error) {
	if len(r.buf) == 0 {
		req, err := r.stream.Recv()
		if err != nil {
			return 0, err
		}
		r.buf = req.GetData()
	}

	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}
