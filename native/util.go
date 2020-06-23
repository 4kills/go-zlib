package native

import "bytes"

const minWritable = 8192

// StreamCloser can indicate whether their underlying stream is closed.
// If so, the StreamCloser must not be used anymore
type StreamCloser interface {
	// IsClosed returns whether the StreamCloser has closed the underlying stream
	IsClosed() bool
}

func startMemAddress(b *bytes.Buffer) *byte {
	if len(b.Bytes()) > 0 {
		return &b.Bytes()[0]
	}
	b.WriteByte(0)
	ptr := &b.Bytes()[0]
	b.Reset()
	return ptr
}
