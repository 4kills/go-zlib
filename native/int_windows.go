// +build windows darwin

package native

// #include <stdint.h>
import "C"

func toInt64(in int64) C.int64_t {
	return C.longlong(in)
}

func intToInt64(in int) C.int64_t {
	return toInt64(int64(in))
}
