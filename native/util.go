package native

/*
#cgo CFLAGS: -I/zlib/
#cgo LDFLAGS: ${SRCDIR}/libs/libz.a

#include "zlib/zlib.h"
*/
import "C"
import "fmt"

const minWritable = 8192
const assumedCompressionFactor = 7

// StreamCloser can indicate whether their underlying stream is closed.
// If so, the StreamCloser must not be used anymore
type StreamCloser interface {
	// IsClosed returns whether the StreamCloser has closed the underlying stream
	IsClosed() bool
}

func grow(b []byte, n int) []byte {
	if cap(b)-len(b) >= n {
		return b
	}

	new := make([]byte, len(b), len(b)+n)

	// supposedly faster than copy(new, b)
	for i := 0; i < len(b); i++ {
		new[i] = b[i]
	}
	return new
}

func determineError(parent error, errCode C.int) error {
	var err error

	switch errCode {
	case C.Z_OK:
		fallthrough
	case C.Z_STREAM_END:
		fallthrough
	case C.Z_NEED_DICT:
		return nil
	case C.Z_STREAM_ERROR:
		err = errStream
	case C.Z_DATA_ERROR:
		err = errData
	case C.Z_MEM_ERROR:
		err = errMem
	case C.Z_VERSION_ERROR:
		err = errMem
	default:
		err = errUnknown
	}

	if parent == nil {
		return err
	}
	return fmt.Errorf("%s: %s", parent.Error(), err.Error())
}

func startMemAddress(b []byte) *byte {
	if len(b) > 0 {
		return &b[0]
	}

	b = append(b, 0)
	ptr := &b[0]
	b = b[0:0]

	return ptr
}
