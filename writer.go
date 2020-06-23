package zlib

import (
	"errors"
	"io"
)

const (
	minCompression = 0
	maxCompression = 9
)

// Writer compresses and writes given data to an underlying io.Writer
type Writer struct {
	w     io.Writer
	level int
}

// NewWriter returns a new Writer with the underlying io.Writer to compress to
func NewWriter(w io.Writer) *Writer {
	nw, _ := NewWriterLevel(w, DefaultCompression)
	return nw
}

// NewWriterLevel performs like NewWriter but you may also specify the compression level
func NewWriterLevel(w io.Writer, level int) (*Writer, error) {
	if level != DefaultCompression && level != HuffmanOnly && (level < minCompression || level > maxCompression) {
		return nil, errors.New("zlib: invalid compression level provided")
	}
	return &Writer{w, level}, nil
}

func (zw *Writer) Write(d []byte) (n int, err error) {
	//TODO: implement
	return 0, nil
}

// Close closes the writer by flushing any unwritten data to the underlying writer.
// You should not forget to call this after being done with the writer.
func (zw *Writer) Close() error {
	//TODO: implement
	return nil
}

// Flush writes any unwritten data to the underlying writer
func (zw *Writer) Flush() error {
	//TODO: implement
	return nil
}

// Reset resets the Writer to the state after being initialized with zlib.NewX(..),
// but with the new underlying writer instead.
func (zw *Writer) Reset(w io.Writer) {
	//TODO: implement
}
