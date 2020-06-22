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

void closeCompressor(longint ptr) {
    z_stream* s = (z_stream*) ptr;
    int ok = deflateEnd(s);

    free(s);

    if(ok != Z_OK) errno = 1;
}

signedint compressData(longint ptr, longint inPtr, signedint inSize, longint outPtr, signedint outSize, signedint* hasCompleted, signedint* processed) {
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

// Compressor using an underlying C zlib stream to compress (deflate) data
type Compressor struct {
	ptr          int64
	hasCompleted *int
	processed    *int
}

func (c *Compressor) Write(b []byte) (int, error) {
	return 0, nil
}
