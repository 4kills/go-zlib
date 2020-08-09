#include "zlib.h"
#include <stdlib.h>
#include <stdint.h>

typedef unsigned char b;

z_stream* newStream();

void freeMem(z_stream* s);

int64_t getProcessed(z_stream* s, int64_t inSize);

int64_t getCompressed(z_stream* s, int64_t outSize);

void prepare(z_stream* s,  int64_t inPtr, int64_t inSize, int64_t outPtr, int64_t outSize);