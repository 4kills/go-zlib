package native

/*
#include "util.h"
#include <stdlib.h>

z_stream* newStream() {
	return (z_stream*) calloc(1, sizeof(z_stream));
}

longint getProcessed(z_stream* s, longint inSize) {
	return inSize - s->avail_in;
}

longint getCompressed(z_stream* s, longint outSize) {
	return outSize - s->avail_out;
}

*/
import "C"

type processor struct {
	s            *C.z_stream
	hasCompleted bool
	processed    int
	isClosed     bool
}

func newProcessor() processor {
	return processor{C.newStream(), false, 0, false}
}

func (p *processor) prepare(inPtr uintptr, inSize int, outPtr uintptr, outSize int) {
	C.prepare(
		p.s,
		C.longlong(inPtr),
		C.longlong(inSize),
		C.longlong(outPtr),
		C.longlong(outSize),
	)
}

func (p *processor) updateProcessed(inSize int) {
	p.processed = int(C.getProcessed(p.s, C.longlong(inSize)))
}

func (p *processor) compressed(outSize int) int {
	return int(C.getCompressed(p.s, C.longlong(outSize)))
}
