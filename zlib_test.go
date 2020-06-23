package zlib

import (
	"bytes"
	"testing"
)

func TestShortString(t *testing.T) {
	var b bytes.Buffer
	in := []byte("hello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\n")
	w, _ := NewWriter(&b)
	w.Write(in)
	w.Close()

	r, _ := NewReader(&b)
	out, _ := r.ReadBytes(b.Bytes())
	r.Close()

	if len(out) != len(in) {
		t.Errorf("inequal size: want %d; got %d", len(in), len(out))
	}
}
