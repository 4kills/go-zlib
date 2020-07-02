package native

/*
#cgo CFLAGS: -I/zlib/
#cgo LDFLAGS: ${SRCDIR}/libs/libz.a

#include "zlib.h"

// I have no idea why I have to wrap just this function but otherwise cgo won't compile
int defInit2(z_stream* s, int lvl, int method, int windowBits, int memLevel, int strategy) {
	return deflateInit2(s, lvl, method, windowBits, memLevel, strategy);
}
*/
import "C"
import (
	"fmt"
)

const defaultWindowBits = 15
const defaultMemLevel = 8

// Compressor using an underlying C zlib stream to compress (deflate) data
type Compressor struct {
	p     processor
	level int
}

// IsClosed returns whether the StreamCloser has closed the underlying stream
func (c *Compressor) IsClosed() bool {
	return c.p.isClosed
}

// NewCompressor returns and initializes a new Compressor with zlib compression stream initialized
func NewCompressor(lvl int) (*Compressor, error) {
	return NewCompressorStrategy(lvl, int(C.Z_DEFAULT_STRATEGY))
}

// NewCompressorStrategy returns and initializes a new Compressor with given level and strategy
// with zlib compression stream initialized
func NewCompressorStrategy(lvl, strat int) (*Compressor, error) {
	p := newProcessor()

	if ok := C.defInit2(p.s, C.int(lvl), C.Z_DEFLATED, C.int(defaultWindowBits), C.int(defaultMemLevel), C.int(strat)); ok != C.Z_OK {
		return nil, determineError(fmt.Errorf("%s: %s", errInitialize.Error(), "compression level might be invalid"), ok)
	}

	return &Compressor{p, lvl}, nil
}

// Close closes the underlying zlib stream and frees the allocated memory
func (c *Compressor) Close() error {
	ok := C.deflateEnd(c.p.s)

	c.p.close()

	if ok != C.Z_OK {
		return determineError(errClose, ok)
	}
	return nil
}

// Compress compresses the given data and returns it as byte slice
func (c *Compressor) Compress(in []byte) ([]byte, error) {
	condition := func() bool {
		return !c.p.hasCompleted
	}

	zlibProcess := func() C.int {
		return C.deflate(c.p.s, C.Z_FINISH)
	}

	specificReset := func() C.int {
		return C.deflateReset(c.p.s)
	}

	_, b, err := c.p.process(
		in,
		make([]byte, 0, len(in)/assumedCompressionFactor),
		condition,
		zlibProcess,
		specificReset,
	)
	return b, err
}
