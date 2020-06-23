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

	return nil
}

func (r *reader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, errNoInput
	}

	buf := new(bytes.Buffer)
	buf.Grow(len(p))
	if _, err := io.Copy(buf, r); err != nil {
		return 0, err
	}

	// get compressed slice here and then wrap in buffer and io.Copy

	return len(p), nil
}

func (zr *reader) Reset(r io.Reader) {
	zr.r = r
}

// NewReader returns a new reader, reading from r. It decompresses read data.
func NewReader(r io.Reader) (io.ReadCloser, error) {
	//c, err := native.NewDecompressor()
	c := &native.Decompressor{} //dummy
	return &reader{r, c}, nil
}
