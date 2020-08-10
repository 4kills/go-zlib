package zlib

import (
	"bytes"
	"compress/zlib"
	"io"
	"testing"
)

// UNIT TESTS

func TestRead_SufficientBuffer(t *testing.T) {
	b := &bytes.Buffer{}
	out := &bytes.Buffer{}
	w := zlib.NewWriter(b)

	r, err := zlib.NewReader(b)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	read := func() {
		p := make([]byte, 1e+4)
		n, err := r.Read(p)
		if err != nil && err != io.EOF {
			t.Error(err)
			t.Error(n)
		}
		out.Write(p[:n])
	}

	_, err = w.Write(shortString)
	w.Flush()

	read()

	_, err = w.Write(shortString)
	w.Close()

	read()

	sliceEquals(t, append(shortString, shortString...), out.Bytes())
}
