package zlib

import (
	"bytes"
	"compress/zlib"
	"io"
	"testing"
)

// UNIT TESTS

func TestReadBytes(t *testing.T) {
	b := &bytes.Buffer{}
	w := zlib.NewWriter(b)
	w.Write(longString)
	w.Close()

	r, err := NewReader(nil)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	_, act, err := r.ReadBuffer(b.Bytes(), nil)
	if err != nil {
		t.Error(err)
	}

	sliceEquals(t, longString, act)
}

func initTestRead(t *testing.T, bufferSize int) (*bytes.Buffer, *zlib.Writer, *Reader, func(r *Reader) error) {
	b := &bytes.Buffer{}
	out := &bytes.Buffer{}
	w := zlib.NewWriter(b)

	r, err := NewReader(b)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	read := func(r *Reader) error {
		p := make([]byte, bufferSize)
		n, err := r.Read(p)
		if err != nil && err != io.EOF {
			t.Error(err)
			t.Error(n)
			t.FailNow()
		}
		out.Write(p[:n])
		return err // io.EOF or nil
	}

	return out, w, r, read
}

func TestRead_SufficientBuffer(t *testing.T) {
	out, w, r, read := initTestRead(t, 1e+4)
	defer r.Close()

	w.Write(shortString)
	w.Flush()

	read(r)

	w.Write(shortString)
	w.Close()

	read(r)

	sliceEquals(t, append(shortString, shortString...), out.Bytes())
}

func TestRead_SmallBuffer(t *testing.T) {
	out, w, r, read := initTestRead(t, 1)
	defer r.Close()

	w.Write(shortString)
	w.Write(shortString)
	w.Close()

	for {
		err := read(r)
		if err == io.EOF {
			break
		}
	}

	sliceEquals(t, append(shortString, shortString...), out.Bytes())
}
