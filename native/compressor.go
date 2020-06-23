package native

/*
#cgo CFLAGS: -I/zlib/
#cgo LDFLAGS: libs/libz.a

#include "util.h"
#include <stdlib.h>

// Returns ptr to stream
longint initCompressor(signedint level) {
    z_stream* s = (z_stream*) calloc(1, sizeof(z_stream));
    int ok = deflateInit(s, level);

    if (ok != Z_OK) {
        errno = 1;
        return -1;
    }

    return (longint) s;
}

void resetCompressor(longint ptr) {
    z_stream* s = (z_stream*) ptr;
    int ok = deflateReset(s);

    if (ok != Z_OK) {
        errno = 1;
    }
}

void closeCompressor(longint ptr) {
    z_stream* s = (z_stream*) ptr;
    int ok = deflateEnd(s);

    free(s);

    if(ok != Z_OK) errno = 1;
}

signedint compressData(longint ptr, longint inPtr, longint inSize, longint outPtr, longint outSize, signedint* hasCompleted, longint* processed) {
    z_stream* s = (z_stream*) ptr;

    prepare(s, inPtr, inSize, outPtr, outSize);

    int ok = deflate (s, Z_FINISH);

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
*/
import "C"
import (
	"bytes"
	"fmt"
	"unsafe"
)

// Compressor using an underlying C zlib stream to compress (deflate) data
type Compressor struct {
	ptr          int64
	hasCompleted int
	processed    int64
	level        int
	isClosed     bool
}

// IsClosed returns whether the StreamCloser has closed the underlying stream
func (c *Compressor) IsClosed() bool {
	return c.isClosed
}

// NewCompressor returns and initializes a new Compressor with zlib compression stream initialized
func NewCompressor(lvl int) (*Compressor, error) {
	ptr, err := C.initCompressor(C.int(lvl))
	if err != nil {
		C.clearError()
		return nil, fmt.Errorf(errInitialize.Error(), ": compression level might be invalid")
	}

	return &Compressor{int64(ptr), 0, 0, lvl, false}, nil
}

// Close closes the underlying zlib stream and frees the allocated memory
func (c *Compressor) Close() error {
	c.isClosed = true
	_, err := C.closeCompressor(C.longlong(c.ptr))
	if err != nil {
		C.clearError()
		return errClose
	}
	return nil
}

// Compress compresses the given data and returns it as byte slice
func (c *Compressor) Compress(in []byte) ([]byte, error) {
	inMem := &in[0]
	inIdx := 0

	var buf bytes.Buffer
	outIdx := 0

	for c.hasCompleted == 0 {
		buf.Grow(minWritable)

		outMem := startMemAddress(&buf)

		readMem := uintptr(unsafe.Pointer(inMem)) + uintptr(inIdx)
		writeMem := uintptr(unsafe.Pointer(outMem)) + uintptr(outIdx)
		compressed, err := C.compressData(C.longlong(c.ptr), C.longlong(readMem), C.longlong(len(in)-inIdx), C.longlong(writeMem), C.longlong(buf.Cap()-outIdx), (*C.int)(unsafe.Pointer(&c.hasCompleted)), (*C.longlong)(unsafe.Pointer(&c.processed)))
		if err != nil {
			C.clearError()
			return nil, errProcess
		}

		inIdx += int(c.processed)
		outIdx += int(compressed)
	}

	c.processed = 0
	c.hasCompleted = 0

	_, err := C.resetCompressor(C.longlong(c.ptr))
	if err != nil {
		C.clearError()
		return buf.Bytes(), errReset
	}

	return buf.Bytes(), nil
}
