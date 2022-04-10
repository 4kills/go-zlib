package zlib

import (
	"io"

	"github.com/4kills/go-zlib/native"
)

const (
	minCompression = 0
	maxCompression = 9

	minStrategy = 0
	maxStrategy = 4
)

// Writer compresses and writes given data to an underlying io.Writer
type Writer struct {
	w          io.Writer
	level      int
	strategy   int
	compressor *native.Compressor
}

// NewWriter returns a new Writer with the underlying io.Writer to compress to.
// w may be nil if you only plan on using WriteBuffer.
// Panics if the underlying c stream cannot be allocated which would indicate a severe error
// not only for this library but also for the rest of your code.
func NewWriter(w io.Writer) *Writer {
	zw, err := NewWriterLevel(w, DefaultCompression)
	if err != nil {
		panic(err)
	}
	return zw
}

// NewWriterLevel performs like NewWriter but you may also specify the compression level.
// w may be nil if you only plan on using WriteBuffer.
func NewWriterLevel(w io.Writer, level int) (*Writer, error) {
	return NewWriterLevelStrategy(w, level, DefaultStrategy)
}

// NewWriterLevelDict does exactly like NewWriterLevel as of NOW.
// This will change once custom dicionaries are implemented.
// This function has been added for compatibility with the std lib.
func NewWriterLevelDict(w io.Writer, level int, dict []byte) (*Writer, error) {
	return NewWriterLevel(w, level)
}

// NewWriterLevelStrategyDict does exactly like NewWriterLevelStrategy as of NOW.
// This will change once custom dicionaries are implemented.
// This function has been added mainly for completeness' sake.
func NewWriterLevelStrategyDict(w io.Writer, level, strategy int, dict []byte) (*Writer, error) {
	return NewWriterLevelStrategy(w, level, strategy)
}

// NewWriterLevelStrategy performs like NewWriter but you may also specify the compression level and strategy.
// w may be nil if you only plan on using WriteBuffer.
func NewWriterLevelStrategy(w io.Writer, level, strategy int) (*Writer, error) {
	if level != DefaultCompression && (level < minCompression || level > maxCompression) {
		return nil, errInvalidLevel
	}
	if strategy < minStrategy || strategy > maxStrategy {
		return nil, errInvalidStrategy
	}

	c, err := native.NewCompressorStrategy(level, strategy)

	return &Writer{w, level, strategy, c}, err
}

// WriteBuffer takes uncompressed data in, compresses it to out and returns out sliced accordingly.
// In most cases (if the compressed data is smaller than the uncompressed)
// an out buffer of size len(in) should be sufficient.
// If you pass nil for out, this function will try to allocate a fitting buffer.
// Use this for whole-buffered, in-memory data.
func (zw *Writer) WriteBuffer(in, out []byte) ([]byte, error) {
	if len(in) == 0 {
		return nil, errNoInput
	}
	if err := checkClosed(zw.compressor); err != nil {
		return nil, err
	}

	if out == nil {
		return zw.compressor.Compress(in, make([]byte, len(in)+16))
	}

	return zw.compressor.Compress(in, out)
}

// Write compresses the given data p and writes it to the underlying io.Writer.
// The data is not necessarily written to the underlying writer, if no Flush is called.
// It returns the number of *uncompressed* bytes written to the underlying io.Writer in case of err = nil,
// or the number of *compressed* bytes in case of err != nil.
// Please consider using WriteBuffer as it might be more convenient for your use case.
func (zw *Writer) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return -1, errNoInput
	}
	if err := checkClosed(zw.compressor); err != nil {
		return -1, err
	}

	out, err := zw.compressor.CompressStream(p)
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
	return len(p), nil
}

// Close closes the writer by flushing any unwritten data to the underlying writer.
// You should not forget to call this after being done with the writer.
func (zw *Writer) Close() error {
	if err := checkClosed(zw.compressor); err != nil {
		return err
	}

	b, err := zw.compressor.Close()
	if err != nil {
		return err
	}

	if zw.w == nil {
		return err
	}

	_, err = zw.w.Write(b)
	return err
}

// Flush writes compressed buffered data to the underlying writer.
func (zw *Writer) Flush() error {
	if err := checkClosed(zw.compressor); err != nil {
		return err
	}
	b, _ := zw.compressor.Flush()
	_, err := zw.w.Write(b)
	return err
}

// Reset flushes the buffered data to the current underyling writer,
// resets the Writer to the state of being initialized with zlib.NewX(..),
// but with the new underlying writer instead.
// This will panic if the writer has already been closed, writer could not be reset or could not write to current
// underlying writer.
func (zw *Writer) Reset(w io.Writer) {
	if err := checkClosed(zw.compressor); err != nil {
		panic(err)
	}

	b, err := zw.compressor.Reset()
	if err != nil {
		panic(err)
	}

	if zw.w == nil {
		zw.w = w
		return
	}

	if _, err := zw.w.Write(b); err != nil {
		panic(err)
	}

	zw.w = w
}
