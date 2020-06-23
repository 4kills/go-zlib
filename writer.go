package zlib

import (
	"io"

	"github.com/4kills/zlib/native"
)

const (
	minCompression = 0
	maxCompression = 9
)

// Writer compresses and writes given data to an underlying io.Writer
type Writer struct {
	w          io.Writer
	level      int
	compressor *native.Compressor
}

// NewWriter returns a new Writer with the underlying io.Writer to compress to.
func NewWriter(w io.Writer) (*Writer, error) {
	return NewWriterLevel(w, DefaultCompression)
}

// NewWriterLevel performs like NewWriter but you may also specify the compression level
func NewWriterLevel(w io.Writer, level int) (*Writer, error) {
	if level != DefaultCompression && level != HuffmanOnly && (level < minCompression || level > maxCompression) {
		return nil, errInvalidLevel
	}

	c, err := native.NewCompressor(level)

	return &Writer{w, level, c}, err
}

func (zw *Writer) Write(d []byte) (int, error) {
	if len(d) == 0 {
		return -1, errNoInput
	}
	if err := checkClosed(zw.compressor); err != nil {
		return -1, err
	}

	out, err := zw.compressor.Compress(d)
	if err != nil {
		return 0, err
	}

	n := 0
	for n < len(out) {
		inc, err := zw.w.Write(out)
		if err != nil {
			return n, err
		}
		n += inc
	}
	return len(d), nil
}

// Close closes the writer by flushing any unwritten data to the underlying writer.
// You should not forget to call this after being done with the writer.
func (zw *Writer) Close() error {
	if err := checkClosed(zw.compressor); err != nil {
		return err
	}

	return zw.compressor.Close()
}

// Flush writes any unwritten data to the underlying writer
func (zw *Writer) Flush() error {
	if err := checkClosed(zw.compressor); err != nil {
		return err
	}
	// if write is successful there is no unwritten data
	return nil
}

// Reset resets the Writer to the state of being initialized with zlib.NewX(..),
// but with the new underlying writer instead.
// This will panic if the writer has already been closed
func (zw *Writer) Reset(w io.Writer) {
	if err := checkClosed(zw.compressor); err != nil {
		panic(err)
	}

	zw.w = w
}
