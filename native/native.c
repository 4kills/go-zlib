#include "zlib/zlib.h"
#include <stdint.h>
#include "zlib/zlib.h"

#define longint int64_t
#define signedint int32_t
#define true 1
#define false 0

typedef unsigned char byte; 
typedef byte bool; 

signedint err; 
bool completed; 
signedint processed;

longint initDecompressor() {
    return internalInit(false, 0);
}

longint initCompressor(int level) {
    return internalInit(true, level);
}

// Returns ptr to stream
longint internalInit(bool shouldCompress, int l) {
    z_stream* s = (z_stream*) calloc(1, sizeof(z_stream)); 
    int ok = (shouldCompress) ? deflateInit(s, l) : inflateInit(s); 

    if (ok != Z_OK) {
        err = 1; 
        return -1; 
    } 
    
    err = 0; 
    return (longint) s; 
}

void closeDecompressor(longint ptr) {
    internalClose(ptr, false); 
}

void closeCompressor(longint ptr) {
    internalClose(ptr, true); 
}

void internalClose(longint ptr, bool shouldCompress) {
    z_stream* s = (z_stream*) ptr; 
    int ok = (shouldCompress) ? deflateEnd(s) : inflateEnd(s); 

    free(s); 

    if(ok != Z_OK) err = 1; 
    else err = 0;  
}

signedint compress(longint ptr, longint inPtr, signedint inSize, longint outPtr, signedint outSize) {
    return internalProcess(ptr, inPtr, inSize, outPtr, outSize, true);
}

signedint decompress(longint ptr, longint inPtr, signedint inSize, longint outPtr, signedint outSize) {
    return internalProcess(ptr, inPtr, inSize, outPtr, outSize, false);
}

signedint internalProcess(longint ptr, longint inPtr, signedint inSize, longint outPtr, signedint outSize, bool shouldCompress) {
    z_stream* s = (z_stream*) ptr; 

    s->avail_in = inSize; 
    s->next_in = (byte*) inPtr; 

    s->avail_out = outSize;
    s->next_out = (byte*) outPtr;

    int ok = (shouldCompress) ? deflate (s, Z_FINISH) : inflate(s, Z_PARTIAL_FLUSH); 

    switch (ok)
    {
    case Z_STREAM_END:
        completed = true; 
        break;
    case Z_OK:
        break;
    default:
        err = 1; 
        return -1; 
    }

    err = 0; 
    processed = inSize - s->avail_in; 

    return outSize - s->avail_out; 
}

