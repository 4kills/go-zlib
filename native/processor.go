package native

/*
#include "util.h"
#include <stdlib.h>

z_stream* newStream() {
	return (z_stream*) calloc(1, sizeof(z_stream));
}

void freeMem(z_stream* s) {
	free(s);
}

longint getProcessed(z_stream* s, longint inSize) {
	return inSize - s->avail_in;
}

longint getCompressed(z_stream* s, longint outSize) {
	return outSize - s->avail_out;
}

*/
import "C"
import "unsafe"

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

func (p *processor) close() {
	C.freeMem(p.s)
	p.s = nil
	p.isClosed = true
}

func (p *processor) updateProcessed(inSize int) {
	p.processed = int(C.getProcessed(p.s, C.longlong(inSize)))
}

func (p *processor) compressed(outSize int) int {
	return int(C.getCompressed(p.s, C.longlong(outSize)))
}

func (p *processor) process(in []byte, buf []byte, condition func() bool, zlibProcess func() C.int, specificReset func() C.int) ([]byte, error) {
	inMem := &in[0]
	inIdx := 0

	outIdx := 0

	for condition() {
		buf = grow(buf, minWritable)

		outMem := startMemAddress(buf)

		readMem := uintptr(unsafe.Pointer(inMem)) + uintptr(inIdx)
		readLen := len(in) - inIdx
		writeMem := uintptr(unsafe.Pointer(outMem)) + uintptr(outIdx)
		writeLen := cap(buf) - outIdx

		p.prepare(readMem, readLen, writeMem, writeLen)

		ok := zlibProcess()
		switch ok {
		case C.Z_STREAM_END:
			p.hasCompleted = true
			break
		case C.Z_OK:
			break
		default:
			return nil, errProcess
		}

		p.updateProcessed(readLen)
		compressed := p.compressed(writeLen)

		inIdx += p.processed
		outIdx += int(compressed)
		buf = buf[:outIdx]
	}

	p.processed = 0
	p.hasCompleted = false

	ok := specificReset()
	if ok != C.Z_OK {
		return buf, errReset
	}

	return buf, nil
}
