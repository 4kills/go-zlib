package zlib

import (
	"bytes"
	"compress/zlib"
	"io"
	"testing"
)

// UNIT TESTS

func TestWrite(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewWriter(b)

	_, err := w.Write(shortString)
	if err != nil {
		t.Error(err)
	}
	w.Flush()
	_, err = w.Write(shortString)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	r, err := zlib.NewReader(b)
	if err != nil {
		t.Error(err)
	}

	act := &bytes.Buffer{}
	_, err = io.Copy(act, r)
	if err != nil {
		t.Error(err)
	}

	sliceEquals(t, append(shortString, shortString...), act.Bytes())
}

func TestReset(t *testing.T) {
	b := &bytes.Buffer{}
	w := NewWriter(b)
	defer w.Close()

	_, err := w.Write(shortString)
	if err != nil {
		t.Error(err)
	}
	w.Flush()
	_, err = w.Write(shortString)
	if err != nil {
		t.Error(err)
	}
	w.Reset(b)

	r, err := zlib.NewReader(b)
	if err != nil {
		t.Error(err)
	}

	act := &bytes.Buffer{}
	_, err = io.Copy(act, r)
	if err != nil {
		t.Error(err)
	}

	sliceEquals(t, append(shortString, shortString...), act.Bytes())
}
