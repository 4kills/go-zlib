package native

import (
	"errors"
)

var (
	errClose      = errors.New("native zlib: zlib stream could not be properly closed and freed")
	errInitialize = errors.New("native zlib: zlib stream could not be properly initialized")
	errProcess    = errors.New("native zlib: zlib stream error during in-/deflation")
	errReset      = errors.New("native zlib: zlib stream could not be properly reset")

	errStream  = errors.New("internal state of stream inconsistent: using same stream over mulitiple threads is not advised")
	errData    = errors.New("data corrupted: data not in a suitable format")
	errMem     = errors.New("out of memory")
	errBuf     = errors.New("avail in or avail out zero")
	errVersion = errors.New("inconsistent zlib version")
	errUnknown = errors.New("error code returned by native c functions unknown")

	retry = errors.New("zlib: ")
)
