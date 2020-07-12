package zlib

import (
	"bytes"
	"compress/zlib"
	"io"
	"testing"
)

// real world data benchmarks

const compressedMcPacketsLoc = "https://raw.githubusercontent.com/4kills/zlib_benchmark/master/compressed_mc_packets.json"

var compressedMcPackets [][]byte

func BenchmarkReadBytesAllMcPacketsDefault(b *testing.B) {
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)

	benchmarkReadBytesMcPacketsGeneric(compressedMcPackets, b)
}

func benchmarkReadBytesMcPacketsGeneric(input [][]byte, b *testing.B) {
	r, _ := NewReader(bytes.NewBuffer(compressedMcPackets[0]))
	defer r.Close()

	b.ResetTimer()

	reportBytesPerChunk(input, b)

	for i := 0; i < b.N; i++ {
		for _, v := range input {
			r.ReadBytes(v)
		}
	}
}

func BenchmarkReadAllMcPacketsDefault(b *testing.B) {
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)
	buf := &bytes.Buffer{}
	r, _ := NewReader(buf)

	benchmarkReadMcPacketsGeneric(r, buf, compressedMcPackets, b)
}

func BenchmarkReadAllMcPacketsDefaultStd(b *testing.B) {
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)

	buf := bytes.NewBuffer(compressedMcPackets[0]) // the std library needs this or else I can't create a reader
	r, _ := zlib.NewReader(buf)
	defer r.Close()

	decompressed := make([]byte, 300000)

	b.ResetTimer()

	reportBytesPerChunk(compressedMcPackets, b)

	for i := 0; i < b.N; i++ {
		for _, v := range compressedMcPackets {
			b.StopTimer()
			res, _ := r.(zlib.Resetter)
			res.Reset(bytes.NewBuffer(v), nil) // to make the std reader work
			b.StartTimer()

			r.Read(decompressed)
		}
	}
}

func benchmarkReadMcPacketsGeneric(r io.ReadCloser, underlyingReader *bytes.Buffer, input [][]byte, b *testing.B) {
	defer r.Close()
	out := make([]byte, 300000)

	b.ResetTimer()

	reportBytesPerChunk(input, b)

	for i := 0; i < b.N; i++ {
		for _, v := range input {
			underlyingReader.Write(v)
			r.Read(out)
		}
	}
}

// laboratory condition benchmarks

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

	b.ResetTimer()

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
	buf := &bytes.Buffer{}
	r, _ := NewReader(buf)
	benchmarkReadLevelGeneric(r, buf, input, level, b)
}

func benchmarkReadLevelGeneric(r io.ReadCloser, underlyingReader *bytes.Buffer, input []byte, level int, b *testing.B) {
	w, _ := NewWriterLevel(nil, level)
	defer w.Close()

	compressed, _ := w.WriteBytes(input)

	defer r.Close()

	decompressed := make([]byte, len(input))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		underlyingReader.Write(compressed) // requires some time but only very little compared to the benchmarked method r.Read
		r.Read(decompressed)
	}
}
