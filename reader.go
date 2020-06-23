package zlib

import (
	"bytes"
	"io"

	"github.com/4kills/zlib/native"
)

// Reader decompresses data from an underlying io.Reader or via the ReadBytes method, which should be preferred
type Reader struct {
	r            io.Reader
	decompressor *native.Decompressor
}

// Close closes the Reader by closing and freeing the underlying zlib stream.
// You should not forget to call this after being done with the writer.
func (r *Reader) Close() error {
	if err := checkClosed(r.decompressor); err != nil {
		return err
	}

	return r.decompressor.Close()
}

// ReadBytes takes compressed data p, decompresses it and returns it as new byte slice.
// This method is generally slightly faster than Read.
func (r *Reader) ReadBytes(compressed []byte) ([]byte, error) {
	if len(compressed) == 0 {
		return nil, errNoInput
	}
	if err := checkClosed(r.decompressor); err != nil {
		return nil, err
	}

	return r.decompressor.Decompress(compressed)
}

// Read reads compressed data from the provided Reader into the provided buffer p.
// Please consider using ReadBytes instead, as it is faster
func (r *Reader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, errNoInput
	}
	if err := checkClosed(r.decompressor); err != nil {
		return 0, err
	}

	bufSlice := make([]byte, 0, len(p))
	buf := bytes.NewBuffer(bufSlice)
	if _, err := io.Copy(buf, r.r); err != nil {
		return 0, err
	}

	out, err := r.decompressor.Decompress(buf.Bytes())
	if err != nil {
		return 0, err
	}

	if len(out) <= len(p) {
		copy(p, out)
		return len(out), io.EOF
	}

	copy(p, out[:len(p)])
	return len(p), io.ErrShortBuffer
}

// Reset resets the Reader to the state of being initialized with zlib.NewX(..),
// but with the new underlying reader instead. It allows for reuse of the same reader.
// This will panic if the writer has already been closed
func (r *Reader) Reset(reader io.Reader) {
	if err := checkClosed(r.decompressor); err != nil {
		panic(err)
	}

	r.r = reader
}

// NewReader returns a new reader, reading from r. It decompresses read data.
// r may be nil if you only plan on using ReadBytes
func NewReader(r io.Reader) (*Reader, error) {
	c, err := native.NewDecompressor()
	return &Reader{r, c}, err
}
