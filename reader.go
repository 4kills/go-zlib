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
	inBuffer     *bytes.Buffer
	outBuffer    *bytes.Buffer
	eof          bool
}

// Close closes the Reader by closing and freeing the underlying zlib stream.
// You should not forget to call this after being done with the writer.
func (r *Reader) Close() error {
	if err := checkClosed(r.decompressor); err != nil {
		return err
	}
	r.inBuffer = nil
	r.outBuffer = nil
	return r.decompressor.Close()
}

// ReadBytes takes compressed data p, decompresses it in one go and returns it as new byte slice.
// This method is generally quite faster than Read if you know the output size beforehand.
// If you don't, you can still try to use that method (provide size <= 0) but that might take longer than Read.
// The method also returns the number n of bytes that were processed from the compressed slice.
// If n < len(compressed) and err == nil then only the first n compressed bytes were in
// a suitable zlib format and as such decompressed.
// ReadBytes resets the reader for new decompression.
func (r *Reader) ReadBytes(compressed []byte, size int) (n int, decompressed []byte, err error) {
	if len(compressed) == 0 {
		return 0, nil, errNoInput
	}
	if err := checkClosed(r.decompressor); err != nil {
		return 0, nil, err
	}

	return r.decompressor.Decompress(compressed, size)
}

// Read reads compressed data from the underlying Reader into the provided buffer p.
// To reuse the reader after an EOF condition, you have to Reset it.
// Please consider using ReadBytes for whole-buffered data instead, as it is faster and generally easier to use.
func (r *Reader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, io.ErrShortBuffer
	}
	if err := checkClosed(r.decompressor); err != nil {
		return 0, err
	}
	if r.outBuffer.Len() == 0 && r.eof {
		return 0, io.EOF
	}

	if r.outBuffer.Len() != 0 {
		min := len(p)
		if len(p) < r.outBuffer.Len() {
			min = r.outBuffer.Len()
		}
		copy(p, r.outBuffer.Bytes()[:min])
		r.outBuffer.Next(min)

		var err error
		if r.outBuffer.Len() == 0 && r.eof {
			err = io.EOF
		}
		return min, err
	}

	n, err := r.r.Read(p)
	if err != nil && err != io.EOF {
		return 0, err
	}
	r.inBuffer.Write(p[:n])

	eof, processed, out, err := r.decompressor.DecompressStream(r.inBuffer.Bytes(), p)
	r.eof = eof
	if err != nil {
		return 0, err
	}
	r.inBuffer.Next(processed)

	if r.eof && len(out) <= len(p){
		copy(p, out)
		return len(out), io.EOF
	}

	if len(out) > len(p) {
		copy(p, out[:len(p)])
		r.outBuffer.Write(out[len(p):])
		return len(p), nil
	}

	copy(p, out)
	return len(out), nil
}

// Reset resets the Reader to the state of being initialized with zlib.NewX(..),
// but with the new underlying reader instead. It allows for reuse of the same reader.
// AS OF NOW dict IS NOT USED. It's just there to implement the Resetter interface
// to allow for easy interchangeability with the std lib. Just pass nil.
func (r *Reader) Reset(reader io.Reader, dict []byte) error {
	if err := checkClosed(r.decompressor); err != nil {
		return err
	}

	err := r.decompressor.Reset()

	r.inBuffer = &bytes.Buffer{}
	r.outBuffer = &bytes.Buffer{}
	r.eof = false
	r.r = reader
	return err
}

// NewReader returns a new reader, reading from r. It decompresses read data.
// r may be nil if you only plan on using ReadBytes
func NewReader(r io.Reader) (*Reader, error) {
	c, err := native.NewDecompressor()
	return &Reader{r, c, &bytes.Buffer{}, &bytes.Buffer{}, false}, err
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
