package zlib

import "github.com/4kills/go-zlib/native"

func checkClosed(c native.StreamCloser) error {
	if c.IsClosed() {
		return errIsClosed
	}
	return nil
}
