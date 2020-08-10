package native

/*
#include "zlib.h"

// I have no idea why I have to wrap just this function but otherwise cgo won't compile
int infInit(z_stream* s) {
	return inflateInit(s);
}
*/
import "C"

// Decompressor using an underlying c zlib stream to decompress (inflate) data
type Decompressor struct {
	p processor
}

// IsClosed returns whether the StreamCloser has closed the underlying stream
func (c *Decompressor) IsClosed() bool {
	return c.p.isClosed
}

// NewDecompressor returns and initializes a new Decompressor with zlib compression stream initialized
func NewDecompressor() (*Decompressor, error) {
	p := newProcessor()

	if ok := C.infInit(p.s); ok != C.Z_OK {
		return nil, determineError(errInitialize, ok)
	}

	return &Decompressor{p}, nil
}

// Close closes the underlying zlib stream and frees the allocated memory
func (c *Decompressor) Close() error {
	ok := C.inflateEnd(c.p.s)

	c.p.close()

	if ok != C.Z_OK {
		return determineError(errClose, ok)
	}
	return nil
}

func (c *Decompressor) Reset() error {
	return determineError(errReset, C.inflateReset(c.p.s))
}

func (c *Decompressor) DecompressStream(in, out []byte) (bool, int, []byte, error) {
	hasCompleted := false
	condition := func() bool {
		hasCompleted = c.p.hasCompleted
		return !c.p.hasCompleted && c.p.readable > 0
	}

	zlibProcess := func() C.int {
		return C.inflate(c.p.s, C.Z_SYNC_FLUSH)
	}

	n, b, err := c.p.process(
		in,
		out[0:0],
		condition,
		zlibProcess,
		func() C.int { return 0 },
	)
	return hasCompleted, n, b, err
}

// Decompress decompresses the given data and returns it as byte slice (preferably in one go)
func (c *Decompressor) Decompress(in, out []byte) (int, []byte, error) {
	zlibProcess := func() C.int {
		ok := C.inflate(c.p.s, C.Z_FINISH)
		if ok == C.Z_BUF_ERROR {
			return 10 // retry
		}
		return ok
	}

	specificReset := func() C.int {
		return C.inflateReset(c.p.s)
	}

	if out != nil {
		return c.p.process(
			in,
			out,
			nil,
			zlibProcess,
			specificReset,
		)
	}

	inc := 1
	for {
		n, b, err := c.p.process(
			in,
			make([]byte, 0, len(in)*assumedCompressionFactor*inc),
			nil,
			zlibProcess,
			specificReset,
		)
		if err == retry {
			inc++
			specificReset()
			continue
		}
		return n, b, err
	}
}
