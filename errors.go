package zlib

import (
	"errors"
)

var (
	errIsClosed        = errors.New("zlib: stream is already closed: you may not use this anymore")
	errNoInput         = errors.New("zlib: no input provided: please provide at least 1 element")
	errInvalidLevel    = errors.New("zlib: invalid compression level provided")
	errInvalidStrategy = errors.New("zlib: invalid compression strategy provided")
)
