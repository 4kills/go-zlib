package zlib

import (
	"bytes"
	"compress/zlib"
	"errors"
	"testing"
	"time"
)

// practical benchmarks

func BenchmarkWriteBytesMcPacketsDefault(b *testing.B) {
	compressedMcPackets := loadPackets("test/mc_packets/decompressed_mc_packets.json")
	w, _ := NewWriter(nil)

	for i := 0; i < b.N; i++ {
		w.WriteBytes(compressedMcPackets[i%len(compressedMcPackets)])
	}
}

func BenchmarkWriteMcPacketsDefault(b *testing.B) {
	compressedMcPackets := loadPackets("test/mc_packets/decompressed_mc_packets.json")

	buf := bytes.Buffer{}
	w, _ := NewWriter(&buf)

	t := time.Now()

	for _, v := range compressedMcPackets {
		w.Write(v)
	}
	b.Log(time.Now().Sub(t))

	/*
		for i := 0; i < b.N; i++ {
			w.Write(compressedMcPackets[i%len(compressedMcPackets)])
		}*/
}

func BenchmarkWriteMcPacketsDefaultStd(b *testing.B) {
	compressedMcPackets := loadPackets("test/mc_packets/decompressed_mc_packets.json")

	buf := bytes.Buffer{}
	w := zlib.NewWriter(&buf)

	//t := time.Now()
	for _, v := range compressedMcPackets {
		w.Write(v)
		w.Flush()
	}
	//b.Log(time.Now().Sub(t))
	b.Log(buf.Len())
	/*
		for i := 0; i < b.N; i++ {
			w.Write(compressedMcPackets[i%len(compressedMcPackets)])
		}*/
}

// laboratory condition benchmarks

func BenchmarkWriteBytes64BBestCompression(b *testing.B) {
	benchmarkWriteBytesLevel(xByte(64), BestCompression, b)
}

func BenchmarkWriteBytes8192BBestCompression(b *testing.B) {
	benchmarkWriteBytesLevel(xByte(8192), BestCompression, b)
}

func BenchmarkWriteBytes65536BBestCompression(b *testing.B) {
	benchmarkWriteBytesLevel(xByte(65536), BestCompression, b)
}

func BenchmarkWriteBytes64BBestSpeed(b *testing.B) {
	benchmarkWriteBytesLevel(xByte(64), BestSpeed, b)
}

func BenchmarkWriteBytes8192BBestSpeed(b *testing.B) {
	benchmarkWriteBytesLevel(xByte(8192), BestSpeed, b)
}

func BenchmarkWriteBytes65536BBestSpeed(b *testing.B) {
	benchmarkWriteBytesLevel(xByte(65536), BestSpeed, b)
}

func BenchmarkWriteBytes64BDefault(b *testing.B) {
	benchmarkWriteBytesLevel(xByte(64), DefaultCompression, b)
}

func BenchmarkWriteBytes8192BDefault(b *testing.B) {
	benchmarkWriteBytesLevel(xByte(8192), DefaultCompression, b)
}

func BenchmarkWriteBytes65536BDefault(b *testing.B) {
	benchmarkWriteBytesLevel(xByte(65536), DefaultCompression, b)
}

func benchmarkWriteBytesLevel(input []byte, level int, b *testing.B) {
	w, _ := NewWriterLevel(nil, level)
	defer w.Close()

	for i := 0; i < b.N; i++ {
		w.WriteBytes(input)
	}
}

func BenchmarkWrite64BBestCompression(b *testing.B) {
	benchmarkWriteLevel(xByte(64), BestCompression, b)
}

func BenchmarkWrite8192BBestCompression(b *testing.B) {
	benchmarkWriteLevel(xByte(8192), BestCompression, b)
}

func BenchmarkWrite65536BBestCompression(b *testing.B) {
	benchmarkWriteLevel(xByte(65536), BestCompression, b)
}

func BenchmarkWrite64BBestSpeed(b *testing.B) {
	benchmarkWriteLevel(xByte(64), BestSpeed, b)
}

func BenchmarkWrite8192BBestSpeed(b *testing.B) {
	benchmarkWriteLevel(xByte(8192), BestSpeed, b)
}

func BenchmarkWrite65536BBestSpeed(b *testing.B) {
	benchmarkWriteLevel(xByte(65536), BestSpeed, b)
}

func BenchmarkWrite64BDefault(b *testing.B) {
	benchmarkWriteLevel(xByte(64), DefaultCompression, b)
}

func BenchmarkWrite8192BDefault(b *testing.B) {
	benchmarkWriteLevel(xByte(8192), DefaultCompression, b)
}

func BenchmarkWrite65536BDefault(b *testing.B) {
	benchmarkWriteLevel(xByte(65536), DefaultCompression, b)
}

func benchmarkWriteLevel(input []byte, level int, b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, len(input)))
	w, _ := NewWriterLevel(buf, level)
	defer w.Close()

	for i := 0; i < b.N; i++ {
		w.Write(input)
		buf.Reset() //requires almost no time
	}
}

func BenchmarkWrite64BBestCompressionStd(b *testing.B) {
	benchmarkWriteLevelStd(xByte(64), BestCompression, b)
}

func BenchmarkWrite8192BBestCompressionStd(b *testing.B) {
	benchmarkWriteLevelStd(xByte(8192), BestCompression, b)
}

func BenchmarkWrite65536BBestCompressionStd(b *testing.B) {
	benchmarkWriteLevelStd(xByte(65536), BestCompression, b)
}

func BenchmarkWrite64BBestSpeedStd(b *testing.B) {
	benchmarkWriteLevelStd(xByte(64), BestSpeed, b)
}

func BenchmarkWrite8192BBestSpeedStd(b *testing.B) {
	benchmarkWriteLevelStd(xByte(8192), BestSpeed, b)
}

func BenchmarkWrite65536BBestSpeedStd(b *testing.B) {
	benchmarkWriteLevelStd(xByte(65536), BestSpeed, b)
}

func BenchmarkWrite64BDefaultStd(b *testing.B) {
	benchmarkWriteLevelStd(xByte(64), DefaultCompression, b)
}

func BenchmarkWrite8192BDefaultStd(b *testing.B) {
	benchmarkWriteLevelStd(xByte(8192), DefaultCompression, b)
}

func BenchmarkWrite65536BDefaultStd(b *testing.B) {
	benchmarkWriteLevelStd(xByte(65536), DefaultCompression, b)
}

var wrote int

// std library zlib in comparison
func benchmarkWriteLevelStd(input []byte, level int, b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, len(input)))
	w, _ := zlib.NewWriterLevel(buf, level) // std library zlib
	defer w.Close()

	n := 0
	for i := 0; i < b.N; i++ {
		n, _ = w.Write(input)
		w.Flush()
		buf.Reset() //requires almost no time
	}
	wrote = n
}

func xByte(multOf16 int) []byte {
	_16byte := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}

	if multOf16 == 0 || multOf16%16 != 0 {
		panic(errors.New("multOf16 is not a valid multiple of 16"))
	}

	xByte := make([]byte, multOf16)
	for i := 0; i < multOf16; i += 16 {
		copy(xByte[i:i+16], _16byte)
	}

	return xByte
}
