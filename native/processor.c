#include "processor.h"

z_stream* newStream() {
	return (z_stream*) calloc(1, sizeof(z_stream));
}

void freeMem(z_stream* s) {
	free(s);
}

int64_t getProcessed(z_stream* s, int64_t inSize) {
	return inSize - s->avail_in;
}

int64_t getCompressed(z_stream* s, int64_t outSize) {
	return outSize - s->avail_out;
}

void prepare(z_stream* s,  int64_t inPtr, int64_t inSize, int64_t outPtr, int64_t outSize) {
    s->avail_in = inSize;
    s->next_in = (b*) inPtr;

    s->avail_out = outSize;
    s->next_out = (b*) outPtr;
}