package zlib

import "fmt"

var (
	errIsClosed = fmt.Errorf("zlib: stream is already closed: you may not use this anymore")
)
