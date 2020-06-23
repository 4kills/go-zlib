package zlib

import "compress/flate"

const (
	//NoCompression does not compress given input
	NoCompression = flate.NoCompression
	//BestSpeed is fastest but with lowest compression
	BestSpeed = flate.BestSpeed
	//BestCompression is slowest but with best compression
	BestCompression = flate.BestCompression
	//DefaultCompression is a compromise between BestSpeed and BestCompression.
	//The level might change if algorithms change.
	DefaultCompression = flate.DefaultCompression
	//HuffmanOnly only uses Huffman encoding to compress the given data
	HuffmanOnly = flate.HuffmanOnly
)
