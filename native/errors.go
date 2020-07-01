package native

import "fmt"

var (
	errClose      = fmt.Errorf("native zlib: zlib stream could not be properly closed and freed")
	errInitialize = fmt.Errorf("native zlib: zlib stream could not be properly initialized")
	errProcess    = fmt.Errorf("native zlib: zlib stream error during in-/deflation")
	errReset      = fmt.Errorf("native zlib: zlib stream could not be properly reset")

	errStream  = fmt.Errorf("internal state of stream inconsistent: using same stream over mulitiple threads is not advised")
	errData    = fmt.Errorf("data corrupted: data not in a suitable format")
	errMem     = fmt.Errorf("out of memory")
	errVersion = fmt.Errorf("inconsistent zlib version")
	errUnknown = fmt.Errorf("error code returned by native c functions unknown")
)
