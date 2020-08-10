package zlib

import (
	"bytes"
	"io"
	"testing"
)

const repeatCount = 30

var shortString = []byte("hello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\n")
var longString []byte

// INTEGRATION TESTS

// primarly checks properly working EOF conditions of Read
func TestWrite_ReadReadFull_wShortString(t *testing.T) {
	b := testWrite(shortString, t)

	r, err := NewReader(b)
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

func TestReadioCopy_wShortString(t *testing.T) {
	// primarly checks properly working EOF conditions of Read
	b := testWrite(shortString, t)

	r, err := NewReader(b)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	bufOut := bytes.NewBuffer(make([]byte, 0, len(shortString)))
	// doesn't know how much data to read (other than other tests)
	_, err = io.Copy(bufOut, r)
	if err != nil {
		t.Error(err)
	}

	sliceEquals(t, shortString, bufOut.Bytes())
}

func TestWrite_ReadOneGo_wShortString(t *testing.T) {
	b := testWrite(shortString, t)

	r, err := NewReader(b)
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
	b := testWrite(shortString, t)
	out := testReadBytes(b, t)
	sliceEquals(t, shortString, out)
}

func TestWriteBytes_ReadBytes_wShortString(t *testing.T) {
	b := testWriteBytes(shortString, t)
	out := testReadBytes(bytes.NewBuffer(b), t)
	sliceEquals(t, shortString, out)
}

func TestReaderReset(t *testing.T) {
	b := testWriteBytes(shortString, t)

	r, err := NewReader(bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	out := make([]byte, len(shortString))
	r.Read(out)
	sliceEquals(t, shortString, out)

	out = make([]byte, len(shortString))
	r.Reset(bytes.NewBuffer(b), nil)
	r.Read(out)
	sliceEquals(t, shortString, out)
}

func TestHuffmanOnly(t *testing.T) {
	w, err := NewWriterLevelStrategy(nil, DefaultCompression, HuffmanOnly)
	if err != nil {
		t.Error(err)
	}
	defer w.Close()

	b, err := w.WriteBuffer(shortString, make([]byte, len(shortString)))
	if err != nil {
		t.Error(err)
	}

	out := testReadBytes(bytes.NewBuffer(b), t)
	sliceEquals(t, shortString, out)
}

func TestWrite_Read_Repeated(t *testing.T) {
	testWriteReadRepeated(shortString, t)
}

func TestRead_RepeatedContinuous(t *testing.T) {
	testReadRepeatedContinuous(shortString, t)
}

func TestWrite_ReadBytes_Repeated(t *testing.T) {
	testWriteReadBytesRepeated(shortString, t)
}

func TestRead_RepeatedContinuous_wLongString(t *testing.T) {
	makeLongString()
	testReadRepeatedContinuous(longString, t)
}

func TestWrite_ReadBytes_Repeated_wLongString(t *testing.T) {
	makeLongString()
	testWriteReadBytesRepeated(longString, t)
}

func TestWrite_Read_Repeated_wLongString(t *testing.T) {
	makeLongString()
	testWriteReadRepeated(longString, t)
}

func testWriteReadBytesRepeated(input []byte, t *testing.T) {
	rep, b := testWriteRepeated(input, t)

	r, err := NewReader(nil)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	out := bytes.NewBuffer(make([]byte, 0, len(rep)))
	m := b.Len()
	for i := 0; i < repeatCount; i++ {
		_, decomp, err := r.ReadBuffer(b.Next(m / repeatCount), nil)
		if err != nil {
			t.Error(err)
		}
		out.Write(decomp)
	}

	sliceEquals(t, rep, out.Bytes())
}

func testReadRepeatedContinuous(input []byte, t *testing.T) {
	compressed := testWriteBytes(input, t)

	b := bytes.Buffer{}
	r, err := NewReader(&b)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	for i := 0; i < repeatCount; i++ {
		b.Write(compressed)

		out := make([]byte, len(input))
		n, err := r.Read(out)
		if err != nil && err != io.EOF {
			t.Error(err)
		}
		if n != len(input) {
			t.Errorf("read count doesn't match: want %d, got %d", len(input), n)
		}
		err = r.Reset(r.r, nil)
		if err != nil {
			t.Error(err)
		}

		sliceEquals(t, input, out)
	}
}

func testWriteReadRepeated(input []byte, t *testing.T) {
	rep, b := testWriteRepeated(input, t)

	stream := &bytes.Buffer{}
	r, err := NewReader(stream)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	out := bytes.NewBuffer(make([]byte, 0, len(rep)))
	for i := 0; i < repeatCount; i++ {
		o := make([]byte, len(rep)/repeatCount)
		var err error
		i := 0
		for err != io.EOF {
			stream.Write(b.Bytes()[i*b.Len()/repeatCount : (i+1) *b.Len()/repeatCount])
			n := 0
			n, err = r.Read(o)
			out.Write(o[:n])
			if err != nil && err != io.EOF {
				t.Error(err)
				t.FailNow()
			}
		}
		err = r.Reset(r.r, nil)
		if err != nil {
			t.Error(err)
		}
	}

	sliceEquals(t, rep, out.Bytes())
}

func testWriteRepeated(input []byte, t *testing.T) ([]byte, *bytes.Buffer) {
	rep := make([]byte, 0, len(input)*repeatCount)
	for i := 0; i < repeatCount; i++ {
		rep = append(rep, input...)
	}

	var b bytes.Buffer
	w := NewWriter(&b)
	defer w.Close()

	for i := 0; i < repeatCount; i++ {
		_, err := w.Write(input)
		if err != nil {
			t.Error(err)
		}
		w.Reset(w.w)
	}
	return rep, &b
}

func testWriteBytes(input []byte, t *testing.T) []byte {
	w := NewWriter(nil)
	defer w.Close()

	b, err := w.WriteBuffer(input, make([]byte, len(input)))
	if err != nil {
		t.Error(err)
	}
	return b
}

func testReadBytes(b *bytes.Buffer, t *testing.T) []byte {
	r, err := NewReader(nil)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	_, out, err := r.ReadBuffer(b.Bytes(), nil)
	if err != nil {
		t.Error(err)
	}
	return out
}

func testWrite(input []byte, t *testing.T) *bytes.Buffer {
	var b bytes.Buffer
	w := NewWriter(&b)
	defer w.Close()

	_, err := w.Write(input)
	if err != nil {
		t.Error(err)
	}
	return &b
}

// HELPER

func makeLongString() {
	if longString != nil {
		return
	}

	for i := 0; i < 150; i++ {
		longString = append(longString, shortString...)
	}
}

func sliceEquals(t *testing.T, expected, actual []byte) {
	if len(expected) != len(actual) {
		t.Errorf("inequal size: want %d; got %d", len(expected), len(actual))
		return
	}
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("slices differ at index %d: want %d; got %d", i, v, actual[i])
			t.FailNow()
		}
	}
}
