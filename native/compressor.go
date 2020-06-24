package native

/*
#cgo CFLAGS: -I/zlib/
#cgo LDFLAGS: ${SRCDIR}/libs/libz.a

#include "zlib/zlib.h"
#include "util.h"
#include <stdlib.h>

// I have no idea why I have to wrap just this function but otherwise cgo won't compile
int defInit(z_stream* s, int lvl) {
	return deflateInit(s, lvl);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

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
	p := newProcessor()

	ok := C.defInit(p.s, C.int(lvl))
	if ok != C.Z_OK {
		return nil, fmt.Errorf(errInitialize.Error(), ": compression level might be invalid")
	}

	return &Compressor{p, lvl}, nil
}

// Close closes the underlying zlib stream and frees the allocated memory
func (c *Compressor) Close() error {
	ok := C.deflateEnd(c.p.s)

	c.p.close()

	if ok != C.Z_OK {
		return errClose
	}
	return nil
}

// Compress compresses the given data and returns it as byte slice
func (c *Compressor) Compress(in []byte) ([]byte, error) {
	inMem := &in[0]
	inIdx := 0

	outIdx := 0

	buf := make([]byte, 0, len(in)/assumedCompressionFactor)

	for !c.p.hasCompleted {
		buf = grow(buf, minWritable)

		outMem := startMemAddress(buf)

		readMem := uintptr(unsafe.Pointer(inMem)) + uintptr(inIdx)
		readLen := len(in) - inIdx
		writeMem := uintptr(unsafe.Pointer(outMem)) + uintptr(outIdx)
		writeLen := cap(buf) - outIdx

		c.p.prepare(readMem, readLen, writeMem, writeLen)

		ok := C.deflate(c.p.s, C.Z_FINISH)
		switch ok {
		case C.Z_STREAM_END:
			c.p.hasCompleted = true
			break
		case C.Z_OK:
			break
		default:
			return nil, errProcess
		}

		c.p.updateProcessed(readLen)
		compressed := c.p.compressed(writeLen)

		inIdx += c.p.processed
		outIdx += int(compressed)
		buf = buf[:outIdx]
	}

	c.p.processed = 0
	c.p.hasCompleted = false

	ok := C.deflateReset(c.p.s)
	if ok != C.Z_OK {
		return buf, errReset
	}

	return buf, nil
}
