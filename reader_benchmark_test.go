package zlib

import (
	"bytes"
	"compress/zlib"
	"testing"
)

func BenchmarkReadBytes64BBestCompression(b *testing.B) {
	benchmarkReadBytesLevel(xByte(64), BestCompression, b)
}

func BenchmarkReadBytes8192BBestCompression(b *testing.B) {
	benchmarkReadBytesLevel(xByte(8192), BestCompression, b)
}

func BenchmarkReadBytes65536BBestCompression(b *testing.B) {
	benchmarkReadBytesLevel(xByte(65536), BestCompression, b)
}

func BenchmarkReadBytes64BBestSpeed(b *testing.B) {
	benchmarkReadBytesLevel(xByte(64), BestSpeed, b)
}

func BenchmarkReadBytes8192BBestSpeed(b *testing.B) {
	benchmarkReadBytesLevel(xByte(8192), BestSpeed, b)
}

func BenchmarkReadBytes65536BBestSpeed(b *testing.B) {
	benchmarkReadBytesLevel(xByte(65536), BestSpeed, b)
}

func BenchmarkReadBytes64BDefault(b *testing.B) {
	benchmarkReadBytesLevel(xByte(64), DefaultCompression, b)
}

func BenchmarkReadBytes8192BDefault(b *testing.B) {
	benchmarkReadBytesLevel(xByte(8192), DefaultCompression, b)
}

func BenchmarkReadBytes65536BDefault(b *testing.B) {
	benchmarkReadBytesLevel(xByte(65536), DefaultCompression, b)
}

func benchmarkReadBytesLevel(input []byte, level int, b *testing.B) {
	w, _ := NewWriterLevel(nil, level)
	defer w.Close()

	compressed, _ := w.WriteBytes(input)

	r, _ := NewReader(nil)
	defer r.Close()

	for i := 0; i < b.N; i++ {
		r.ReadBytes(compressed)
	}
}

func BenchmarkRead64BBestCompression(b *testing.B) {
	benchmarkReadLevel(xByte(64), BestCompression, b)
}

func BenchmarkRead8192BBestCompression(b *testing.B) {
	benchmarkReadLevel(xByte(8192), BestCompression, b)
}

func BenchmarkRead65536BBestCompression(b *testing.B) {
	benchmarkReadLevel(xByte(65536), BestCompression, b)
}

func BenchmarkRead64BBestSpeed(b *testing.B) {
	benchmarkReadLevel(xByte(64), BestSpeed, b)
}

func BenchmarkRead8192BBestSpeed(b *testing.B) {
	benchmarkReadLevel(xByte(8192), BestSpeed, b)
}

func BenchmarkRead65536BBestSpeed(b *testing.B) {
	benchmarkReadLevel(xByte(65536), BestSpeed, b)
}

func BenchmarkRead64BDefault(b *testing.B) {
	benchmarkReadLevel(xByte(64), DefaultCompression, b)
}

func BenchmarkRead8192BDefault(b *testing.B) {
	benchmarkReadLevel(xByte(8192), DefaultCompression, b)
}

func BenchmarkRead65536BDefault(b *testing.B) {
	benchmarkReadLevel(xByte(65536), DefaultCompression, b)
}

func benchmarkReadLevel(input []byte, level int, b *testing.B) {
	w, _ := NewWriterLevel(nil, level)
	defer w.Close()

	compressed, _ := w.WriteBytes(input)

	r, _ := NewReader(bytes.NewBuffer(compressed))
	defer r.Close()

	decompressed := make([]byte, len(input))

	for i := 0; i < b.N; i++ {
		r.Read(decompressed)
		r.Reset(bytes.NewBuffer(compressed)) // requires some time but only very little compared to the benchmarked method r.Read
	}
}

func BenchmarkRead64BBestCompressionStd(b *testing.B) {
	benchmarkReadLevelStd(xByte(64), BestCompression, b)
}

func BenchmarkRead8192BBestCompressionStd(b *testing.B) {
	benchmarkReadLevelStd(xByte(8192), BestCompression, b)
}

func BenchmarkRead65536BBestCompressionStd(b *testing.B) {
	benchmarkReadLevelStd(xByte(65536), BestCompression, b)
}

func BenchmarkRead64BBestSpeedStd(b *testing.B) {
	benchmarkReadLevelStd(xByte(64), BestSpeed, b)
}

func BenchmarkRead8192BBestSpeedStd(b *testing.B) {
	benchmarkReadLevelStd(xByte(8192), BestSpeed, b)
}

func BenchmarkRead65536BBestSpeedStd(b *testing.B) {
	benchmarkReadLevelStd(xByte(65536), BestSpeed, b)
}

func BenchmarkRead64BDefaultStd(b *testing.B) {
	benchmarkReadLevelStd(xByte(64), DefaultCompression, b)
}

func BenchmarkRead8192BDefaultStd(b *testing.B) {
	benchmarkReadLevelStd(xByte(8192), DefaultCompression, b)
}

func BenchmarkRead65536BDefaultStd(b *testing.B) {
	benchmarkReadLevelStd(xByte(65536), DefaultCompression, b)
}

var read int

func benchmarkReadLevelStd(input []byte, level int, b *testing.B) {
	w, _ := NewWriterLevel(nil, level)
	defer w.Close()

	compressed, _ := w.WriteBytes(input)

	buf := bytes.NewBuffer(compressed)
	r, _ := zlib.NewReader(buf)
	defer r.Close()

	decompressed := make([]byte, len(input))

	n := 0
	for i := 0; i < b.N; i++ {
		n, _ := r.Read(decompressed)
		buf.Write(compressed) // requires some time but only very little compared to the benchmarked method r.Read
	}
	read = n
}
