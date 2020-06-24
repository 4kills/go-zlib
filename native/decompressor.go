package native

/*
#cgo CFLAGS: -I/zlib/
#cgo LDFLAGS: ${SRCDIR}/libs/libz.a

#include "util.h"
#include <stdlib.h>

// Returns ptr to stream
longint initDecompressor() {
    z_stream* s = (z_stream*) calloc(1, sizeof(z_stream));
    int ok = inflateInit(s);

    if (ok != Z_OK) {
		errno = 1;
        return -1;
    }

    return (longint) s;
}

void resetDecompressor(longint ptr) {
    z_stream* s = (z_stream*) ptr;
    int ok = inflateReset(s);

    if (ok != Z_OK) {
        errno = 1;
    }
}

void closeDecompressor(longint ptr) {
    z_stream* s = (z_stream*) ptr;
    int ok = inflateEnd(s);

    free(s);

    if(ok != Z_OK) errno = 1;
}

signedint decompressData(longint ptr, longint inPtr, longint inSize, longint outPtr, longint outSize, signedint* hasCompleted, longint* processed) {
    z_stream* s = (z_stream*) ptr;

    prepare(s, inPtr, inSize, outPtr, outSize);

    int ok = inflate(s, Z_PARTIAL_FLUSH);

    switch (ok) {
    case Z_STREAM_END:
        *hasCompleted = true;
        break;
    case Z_OK:
        break;
    default:
        errno = 1;
        return -1;
    }

    *processed = inSize - s->avail_in;

    return outSize - s->avail_out;
}

import "C"
import (
	"unsafe"
)

// Decompressor using an underlying c zlib stream to decompress (inflate) data
type Decompressor struct {
	ptr          int64
	hasCompleted int
	processed    int64
	isClosed     bool
}

// IsClosed returns whether the StreamCloser has closed the underlying stream
func (c *Decompressor) IsClosed() bool {
	return c.isClosed
}

// NewDecompressor returns and initializes a new Decompressor with zlib compression stream initialized
func NewDecompressor() (*Decompressor, error) {
	ptr, err := C.initDecompressor()
	if err != nil {
		C.clearError()
		return nil, errInitialize
	}

	return &Decompressor{int64(ptr), 0, 0, false}, nil
}

// Close closes the underlying zlib stream and frees the allocated memory
func (c *Decompressor) Close() error {
	c.isClosed = true
	_, err := C.closeDecompressor(C.longlong(c.ptr))
	if err != nil {
		C.clearError()
		return errClose
	}
	return nil
}

// Decompress decompresses the given data and returns it as byte slice
func (c *Decompressor) Decompress(in []byte) ([]byte, error) {
	inMem := &in[0]
	inIdx := 0

	outIdx := 0

	buf := make([]byte, len(in)*assumedCompressionFactor)

	for c.hasCompleted == 0 && len(in)-inIdx > 0 {
		buf = grow(buf, minWritable)

		outMem := startMemAddress(buf)

		readMem := uintptr(unsafe.Pointer(inMem)) + uintptr(inIdx)
		writeMem := uintptr(unsafe.Pointer(outMem)) + uintptr(outIdx)

		compressed, err := C.decompressData(
			C.longlong(c.ptr),
			C.longlong(readMem),
			C.longlong(len(in)-inIdx),
			C.longlong(writeMem),
			C.longlong(cap(buf)-outIdx),
			(*C.int)(unsafe.Pointer(&c.hasCompleted)),
			(*C.longlong)(unsafe.Pointer(&c.processed)),
		)

		if err != nil {
			C.clearError()
			return nil, errProcess
		}

		inIdx += int(c.processed)
		outIdx += int(compressed)
		buf = buf[:outIdx]
	}

	c.processed = 0
	c.hasCompleted = 0

	_, err := C.resetDecompressor(C.longlong(c.ptr))
	if err != nil {
		C.clearError()
		return buf, errReset
	}

	return buf, nil
}*/
