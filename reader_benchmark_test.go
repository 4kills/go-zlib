package zlib

import (
	"bytes"
	"compress/zlib"
	"io"
	"testing"
)

// real world data benchmarks

const compressedMcPacketsLoc = "test/mc_packets/compressed_mc_packets.json"

var compressedMcPackets [][]byte

func BenchmarkReadBytesAllMcPacketsDefault(b *testing.B) {
	b.StopTimer()
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)

	benchmarkReadBytesMcPacketsGeneric(compressedMcPackets, b)
}

func benchmarkReadBytesMcPacketsGeneric(input [][]byte, b *testing.B) {
	r, _ := NewReader(bytes.NewBuffer(compressedMcPackets[0]))
	defer r.Close()

	reportBytesPerChunk(input, b)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range input {
			r.ReadBytes(v)
		}
	}
}

func BenchmarkReadAllMcPacketsDefault(b *testing.B) {
	b.StopTimer()
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)
	buf := &bytes.Buffer{}
	r, _ := NewReader(buf)

	benchmarkReadMcPacketsGeneric(r, buf, compressedMcPackets, b)
}

func BenchmarkReadAllMcPacketsDefaultStd(b *testing.B) {
	b.StopTimer()
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)
	buf := bytes.NewBuffer(compressedMcPackets[0]) // the std lib loses it's shit if buf is empty
	r, _ := zlib.NewReader(buf)

	benchmarkReadMcPacketsGeneric(r, buf, compressedMcPackets[1:2], b)
}

func benchmarkReadMcPacketsGeneric(r io.ReadCloser, underlyingReader *bytes.Buffer, input [][]byte, b *testing.B) {
	reportBytesPerChunk(input, b)

	defer r.Close()

	out := make([]byte, 300000)

	b.StartTimer()

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
	b.StopTimer()
	w, _ := NewWriterLevel(nil, level)
	defer w.Close()

	compressed, _ := w.WriteBytes(input)

	r, _ := NewReader(nil)
	defer r.Close()
	b.StartTimer()

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

func benchmarkReadLevelStd(input []byte, level int, b *testing.B) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{1})
	r, _ := zlib.NewReader(buf)
	buf.Reset()
	benchmarkReadLevelGeneric(r, buf, input, level, b)
}

func benchmarkReadLevelGeneric(r io.ReadCloser, underlyingReader *bytes.Buffer, input []byte, level int, b *testing.B) {
	b.StopTimer()
	w, _ := NewWriterLevel(nil, level)
	defer w.Close()

	compressed, _ := w.WriteBytes(input)

	defer r.Close()

	decompressed := make([]byte, len(input))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		underlyingReader.Write(compressed) // requires some time but only very little compared to the benchmarked method r.Read
		r.Read(decompressed)
	}
}
