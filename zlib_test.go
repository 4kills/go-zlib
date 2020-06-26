package zlib

import (
	"bytes"
	"io"
	"testing"
)

const repeatCount = 30

var shortString = []byte("hello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\n")

// primarly checks properly working EOF conditions of Read
func TestWrite_ReadReadFull_wShortString(t *testing.T) {
	var b bytes.Buffer
	w, err := NewWriter(&b)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	_, err = w.Write(shortString)
	if err != nil {
		t.Error(err)
	}

	r, err := NewReader(&b)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	out := make([]byte, len(shortString))
	_, err = io.ReadFull(r, out)
	if err != nil {
		t.Error(err)
	}

	sliceEquals(t, shortString, out)
}

// primarly checks properly working EOF conditions of Read
func TestReadioCopy_wShortString(t *testing.T) {
	var b bytes.Buffer
	w, err := NewWriter(&b)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	_, err = w.Write(shortString)
	if err != nil {
		t.Error(err)
	}

	r, err := NewReader(&b)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	bufOut := bytes.NewBuffer(make([]byte, 0, len(shortString)))
	_, err = io.Copy(bufOut, r) // doesn't know how much data to read (other than previous tests)
	if err != nil {
		t.Error(err)
	}

	sliceEquals(t, shortString, bufOut.Bytes())
}

func TestWrite_ReadOneGo_wShortString(t *testing.T) {
	var b bytes.Buffer
	w, err := NewWriter(&b)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	_, err = w.Write(shortString)
	if err != nil {
		t.Error(err)
	}

	r, err := NewReader(&b)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	out := make([]byte, len(shortString))
	n, err := r.Read(out)
	if n != len(shortString) {
		t.Errorf("did not read enough bytes: want %d; got %d", len(shortString), n)
	}
	if err != io.EOF && err != nil {
		t.Error(err)
	}

	sliceEquals(t, shortString, out)
}

func TestWrite_ReadBytes_wShortString(t *testing.T) {
	var b bytes.Buffer
	w, err := NewWriter(&b)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	_, err = w.Write(shortString)
	if err != nil {
		t.Error(err)
	}

	r, err := NewReader(nil)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	_, out, err := r.ReadBytes(b.Bytes())
	if err != nil {
		t.Error(err)
	}

	sliceEquals(t, shortString, out)
}

func TestWriteBytes_ReadBytes_wShortString(t *testing.T) {
	w, err := NewWriter(nil)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	b, err := w.WriteBytes(shortString)
	if err != nil {
		t.Error(err)
	}

	r, err := NewReader(nil)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	_, out, err := r.ReadBytes(b)
	if err != nil {
		t.Error(err)
	}

	sliceEquals(t, shortString, out)
}

func TestWrite_Read_Repeated(t *testing.T) {
	rep := make([]byte, 0, len(shortString)*repeatCount)
	for i := 0; i < repeatCount; i++ {
		rep = append(rep, shortString...)
	}

	var b bytes.Buffer
	w, err := NewWriter(&b)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	for i := 0; i < repeatCount; i++ {
		_, err = w.Write(shortString)
		if err != nil {
			t.Error(err)
		}
	}

	r, err := NewReader(&b)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	out := bytes.NewBuffer(make([]byte, 0, len(rep)))
	for i := 0; i < repeatCount; i++ {
		o := make([]byte, len(rep)/repeatCount)
		n, err := r.Read(o)
		out.Write(o[:n])
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
		}
	}

	sliceEquals(t, rep, out.Bytes())
}

func TestWrite_ReadBytes_Repeated(t *testing.T) {
	rep := make([]byte, 0, len(shortString)*repeatCount)
	for i := 0; i < repeatCount; i++ {
		rep = append(rep, shortString...)
	}

	var b bytes.Buffer
	w, err := NewWriter(&b)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	for i := 0; i < repeatCount; i++ {
		_, err = w.Write(shortString)
		if err != nil {
			t.Error(err)
		}
	}

	r, err := NewReader(nil)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	out := bytes.NewBuffer(make([]byte, 0, len(rep)))
	m := b.Len()
	for i := 0; i < repeatCount; i++ {
		_, decomp, err := r.ReadBytes(b.Next(m / repeatCount))
		if err != nil {
			t.Error(err)
		}
		out.Write(decomp)
	}

	sliceEquals(t, rep, out.Bytes())
}

func TestHuffmanOnly(t *testing.T) {
	// ASSUMES READ AND WRITE WORK PROPERLY
	w, err := NewWriterLevel(nil, HuffmanOnly)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	b, err := w.WriteBytes(shortString)
	if err != nil {
		t.Error(err)
	}

	r, err := NewReader(nil)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	_, out, err := r.ReadBytes(b)
	if err != nil {
		t.Error(err)
	}

	sliceEquals(t, shortString, out)
}

func sliceEquals(t *testing.T, expected, actual []byte) {
	if len(expected) != len(actual) {
		t.Errorf("inequal size: want %d; got %d", len(expected), len(actual))
		return
	}
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("slices differ at index %d: want %d; got %d", i, v, actual[i])
		}
	}
}
