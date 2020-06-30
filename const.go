package zlib

const (
	// Compression Levels

	//NoCompression does not compress given input
	NoCompression = 0
	//BestSpeed is fastest but with lowest compression
	BestSpeed = 1
	//BestCompression is slowest but with best compression
	BestCompression = 9
	//DefaultCompression is a compromise between BestSpeed and BestCompression.
	//The level might change if algorithms change.
	DefaultCompression = -1

	// Compression Strategies

	// Filtered is more effective for small (but not all too many) randomly distributed values.
	// It forces more Huffman encoding. Use it for filtered data. It's between Default and Huffman only.
	Filtered = 1
	//HuffmanOnly only uses Huffman encoding to compress the given data
	HuffmanOnly = 2
	//RLE (run-length encoding) limits match distance to one, thereby being almost as fast as HuffmanOnly
	// but giving better compression for PNG data.
	RLE = 3
	//Fixed disallows dynamic Huffman codes, thereby making it a simpler decoder
	Fixed = 4
	// DefaultStrategy is the default compression strategy that should be used for most appliances
	DefaultStrategy = 0
)
