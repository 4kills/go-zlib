package native

/*
#cgo CFLAGS: -I/zlib/
#cgo LDFLAGS: libs/libz.a

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

void closeDecompressor(longint ptr) {
    z_stream* s = (z_stream*) ptr;
    int ok = inflateEnd(s);

    free(s);

    if(ok != Z_OK) errno = 1;
}

signedint decompressData(longint ptr, longint inPtr, signedint inSize, longint outPtr, signedint outSize) {
    z_stream* s = (z_stream*) ptr;

    prepare(s, inPtr, inSize, outPtr, outSize);

    int ok = inflate(s, Z_PARTIAL_FLUSH);

    switch (ok) {
    case Z_STREAM_END:
        break;
    case Z_OK:
        break;
    default:
        errno = 1;
        return -1;
    }

    return outSize - s->avail_out;
}
*/
import "C"

// Decompressor using an underlying c zlib stream to decompress (inflate) data
type Decompressor struct {
	ptr int64
}
