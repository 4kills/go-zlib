#include "zlib/zlib.h"
#include <stdint.h>
#include <errno.h>

#define longint int64_t
#define signedint int32_t
#define true 1
#define false 0

typedef unsigned char byte;

inline void clearError() {
    errno = 0; 
}

inline void prepare(z_stream* s,  longint inPtr, signedint inSize, longint outPtr, signedint outSize) {
    s->avail_in = inSize;
    s->next_in = (byte*) inPtr;

    s->avail_out = outSize;
    s->next_out = (byte*) outPtr;
}