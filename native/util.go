package native

const minWritable = 8192

// StreamCloser can indicate whether their underlying stream is closed.
// If so, the StreamCloser must not be used anymore
type StreamCloser interface {
	// IsClosed returns whether the StreamCloser has closed the underlying stream
	IsClosed() bool
}

func grow(b []byte, n int) []byte {
	if cap(b)-len(b) >= n {
		return b
	}

	new := make([]byte, len(b), len(b)+n)

	for i := 0; i < len(b); i++ {
		new[i] = (b)[i]
	}
	return new
}

func startMemAddress(b []byte) *byte {
	if len(b) > 0 {
		return &b[0]
	}
	b = append(b, 0)
	ptr := &b[0]
	b = b[0:0]
	return ptr
}
