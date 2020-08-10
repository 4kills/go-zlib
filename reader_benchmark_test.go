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
	loadPacketsIfNil(&decompressedMcPackets, decompressedMcPacketsLoc)

	benchmarkReadBytesMcPacketsGeneric(compressedMcPackets, b)
}

func benchmarkReadBytesMcPacketsGeneric(input [][]byte, b *testing.B) {
	r, _ := NewReader(nil)
	defer r.Close()

	b.ResetTimer()

	reportBytesPerChunk(input, b)

	for i := 0; i < b.N; i++ {
		for j, v := range input {
			r.ReadBuffer(v, make([]byte, len(decompressedMcPackets[j])))
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
			res.Reset(bytes.NewBuffer(v), nil)
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
			b.StopTimer()
			res, _ := r.(Resetter)
			res.Reset(bytes.NewBuffer(v), nil)
			b.StartTimer()

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

	compressed, _ := w.WriteBuffer(input, make([]byte, len(input)))

	r, _ := NewReader(nil)
	defer r.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.ReadBuffer(compressed, nil)
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

	compressed, _ := w.WriteBuffer(input, make([]byte, len(input)))

	defer r.Close()

	decompressed := make([]byte, len(input))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		res := r.(Resetter)
		res.Reset(bytes.NewBuffer(compressed), nil)
		b.StartTimer()
		var err error
		for err != io.EOF {
			_, err = r.Read(decompressed)
		}
	}
}
