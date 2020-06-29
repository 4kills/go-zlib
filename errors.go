package zlib

import (
	"fmt"
)

var (
	errIsClosed        = fmt.Errorf("zlib: stream is already closed: you may not use this anymore")
	errNoInput         = fmt.Errorf("zlib: no input provided: please provide at least 1 element")
	errInvalidLevel    = fmt.Errorf("zlib: invalid compression level provided")
	errInvalidStrategy = fmt.Errorf("zlib: invalid compression strategy provided")
)
