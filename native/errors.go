package native

import "fmt"

var (
	errClose      = fmt.Errorf("native zlib: zlib stream could not be properly closed and freed")
	errInitialize = fmt.Errorf("native zlib: zlib stream could not be properly initialized")
	errIsClosed   = fmt.Errorf("native zlib: attempted to use closed stream")
	errProcess    = fmt.Errorf("native zlib: zlib stream error during in-/deflation")
	errReset      = fmt.Errorf("native zlib: zlib stream could not be properly reset")
)
