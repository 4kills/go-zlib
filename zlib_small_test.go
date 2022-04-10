package zlib

import (
	"bytes"
	"log"
	"testing"
)

var tinyString = []byte("wowie")
var shorterString = []byte("ajcpr83;r;729l3yfgn260mpod2zg7z9p")

func TestWriteBytes_ReadBytes_wTinyString(t *testing.T) {
	b := testWriteBytes(tinyString, t)
	out := testReadBytes(bytes.NewBuffer(b), t)
	sliceEquals(t, tinyString, out)
}
func TestWriteBytes_ReadBytes_wShorterString(t *testing.T) {
	b := testWriteBytes(shorterString, t)
	log.Println(b)
	out := testReadBytes(bytes.NewBuffer(b), t)
	sliceEquals(t, shorterString, out)
}
