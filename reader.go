package zlib

import (
	"bytes"
	"io"

	"github.com/4kills/zlib/native"
)

type reader struct {
	r            io.Reader
	decompressor *native.Decompressor
}

func (r *reader) Close() error {
	if err := checkClosed(r.decompressor); err != nil {
		return err
	}

	return r.decompressor.Close()
}

func (r *reader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, errNoInput
	}
	if err := checkClosed(r.decompressor); err != nil {
		return 0, err
	}

	buf := new(bytes.Buffer)
	buf.Grow(len(p))
	if _, err := io.Copy(buf, r); err != nil {
		return 0, err
	}

	out, err := r.decompressor.Decompress(buf.Bytes())
	if len(out) <= len(p) {
		p = out
		return len(out), err
	}
	p = out[:len(p)]
	return len(p), nil
}

// Reset resets the Reader to the state of being initialized with zlib.NewX(..),
// but with the new underlying reader instead. It allows for reuse of the same reader.
// This will panic if the writer has already been closed
func (r *reader) Reset(reader io.Reader) {
	if err := checkClosed(r.decompressor); err != nil {
		panic(err)
	}

	r.r = reader
}

// NewReader returns a new reader, reading from r. It decompresses read data.
func NewReader(r io.Reader) (io.ReadCloser, error) {
	c, err := native.NewDecompressor()
	return &reader{r, c}, err
}
