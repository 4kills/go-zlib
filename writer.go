package zlib

import (
	"io"

	"github.com/4kills/zlib/native"
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
// w may be nil if you only plan on using WriteBytes.
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
// w may be nil if you only plan on using WriteBytes.
func NewWriterLevel(w io.Writer, level int) (*Writer, error) {
	return NewWriterLevelStrategy(w, level, DefaultStrategy)
}

// NewWriterLevelDict does exactly like NewWriterLevel as of NOW.
// This will change once custom dicionaries are implemented.
// This function has been added for compatability with the std lib.
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
// w may be nil if you only plan on using WriteBytes.
func NewWriterLevelStrategy(w io.Writer, level, strategy int) (*Writer, error) {
	if level != DefaultCompression && (level < minCompression || level > maxCompression) {
		return nil, errInvalidLevel
	}
	if strategy < minCompression || strategy > maxCompression {
		return nil, errInvalidStrategy
	}

	c, err := native.NewCompressorStrategy(level, strategy)

	return &Writer{w, level, strategy, c}, err
}

// WriteBytes takes uncompressed data p, compresses it and returns it as new byte slice.
func (zw *Writer) WriteBytes(p []byte) ([]byte, error) {
	if len(p) == 0 {
		return nil, errNoInput
	}
	if err := checkClosed(zw.compressor); err != nil {
		return nil, err
	}

	return zw.compressor.Compress(p)
}

// Write compresses the given data p and writes it to the underlying io.Writer.
// It returns the number of *compressed* bytes written to the underlying io.Writer.
// Please consider using WriteBytes as it might be more convenient for your use case.
func (zw *Writer) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return -1, errNoInput
	}
	if err := checkClosed(zw.compressor); err != nil {
		return -1, err
	}

	out, err := zw.compressor.Compress(p)
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
	return len(out), nil
}

// Close closes the writer by flushing any unwritten data to the underlying writer.
// You should not forget to call this after being done with the writer.
func (zw *Writer) Close() error {
	if err := checkClosed(zw.compressor); err != nil {
		return err
	}

	return zw.compressor.Close()
}

// Flush does NOTHING. Write always flushes content to the underlying writer.
// This method just exists for easy interchangability with the go std zlib library
func (zw *Writer) Flush() error {
	if err := checkClosed(zw.compressor); err != nil {
		return err
	}
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
