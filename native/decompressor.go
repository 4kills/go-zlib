package native

/*
#cgo CFLAGS: -I/zlib/
#cgo LDFLAGS: ${SRCDIR}/libs/libz.a

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

// Decompress decompresses the given data and returns it as byte slice
func (c *Decompressor) Decompress(in []byte) (int, []byte, error) {
	condition := func() bool {
		return !c.p.hasCompleted && c.p.readable > 0
	}

	zlibProcess := func() C.int {
		return C.inflate(c.p.s, C.Z_PARTIAL_FLUSH)
	}

	specificReset := func() C.int {
		return C.inflateReset(c.p.s)
	}

	return c.p.process(
		in,
		make([]byte, 0, len(in)*assumedCompressionFactor),
		condition,
		zlibProcess,
		specificReset,
	)
}
