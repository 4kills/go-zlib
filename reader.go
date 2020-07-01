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
	buffer       *bytes.Buffer
}

// Close closes the Reader by closing and freeing the underlying zlib stream.
// You should not forget to call this after being done with the writer.
func (r *Reader) Close() error {
	if err := checkClosed(r.decompressor); err != nil {
		return err
	}
	r.buffer = nil
	return r.decompressor.Close()
}

// ReadBytes takes compressed data p, decompresses it and returns it as new byte slice.
// It also returns the number n of bytes that were processed from the compressed slice.
// If n < len(compressed) and err == nil then only the first n compressed bytes were in
// a suitable zlib format and as such decompressed.
// This method is generally slightly faster than Read.
func (r *Reader) ReadBytes(compressed []byte) (n int, decompressed []byte, err error) {
	if len(compressed) == 0 {
		return 0, nil, errNoInput
	}
	if err := checkClosed(r.decompressor); err != nil {
		return 0, nil, err
	}

	return r.decompressor.Decompress(compressed)
}

// Read reads compressed data from the provided Reader into the provided buffer p.
// Please consider using ReadBytes instead, as it is slightly faster and generally easier to use
func (r *Reader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, errNoInput
	}
	if err := checkClosed(r.decompressor); err != nil {
		return 0, err
	}

	r.buffer.Grow(len(p))
	if _, err := io.Copy(r.buffer, r.r); err != nil {
		return 0, err
	}

	if r.buffer.Len() == 0 {
		return 0, io.EOF
	}

	processed, out, err := r.decompressor.Decompress(r.buffer.Bytes())
	if err != nil {
		return 0, err
	}
	r.buffer.Next(processed)

	if len(out) > len(p) {
		copy(p, out[:len(p)])
		return len(p), io.ErrShortBuffer
	}

	copy(p, out)
	if r.buffer.Len() == 0 {
		return len(out), io.EOF
	}
	return len(out), nil
}

// Reset resets the Reader to the state of being initialized with zlib.NewX(..),
// but with the new underlying reader instead. It allows for reuse of the same reader.
// AS OF NOW dict IS NOT USED. It's just there to implement the Resetter interface
// to allow for easy interchangeablility with the std lib. Just pass nil.
func (r *Reader) Reset(reader io.Reader, dict []byte) error {
	if err := checkClosed(r.decompressor); err != nil {
		return err
	}

	r.buffer = &bytes.Buffer{}
	r.r = reader
	return nil
}

// NewReader returns a new reader, reading from r. It decompresses read data.
// r may be nil if you only plan on using ReadBytes
func NewReader(r io.Reader) (*Reader, error) {
	c, err := native.NewDecompressor()
	return &Reader{r, c, &bytes.Buffer{}}, err
}

// NewReaderDict does exactly like NewReader as of NOW.
// This will change once custom dicionaries are implemented.
// This function has been added for compatibility with the std lib.
func NewReaderDict(r io.Reader, dict []byte) (*Reader, error) {
	return NewReader(r)
}

// Resetter resets the zlib.Reader returned by NewReader by assigning a new underyling reader,
// discarding any buffered data from the previous reader.
// This interface is mainly for compatibility with the std lib
type Resetter interface {
	// Reset resets the Reader to the state of being initialized with zlib.NewX(..),
	// but with the new underlying reader and dict instead. It allows for reuse of the same reader.
	Reset(r io.Reader, dict []byte) error
}
