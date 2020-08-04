package native

/* #include <stdint.h>

static int64_t convert(long long integer) {
	return (int64_t) integer;
}
*/
import "C"

func toInt64(in int64) C.int64_t {
	return C.convert(C.longlong(in))
}

func intToInt64(in int) C.int64_t {
	return toInt64(int64(in))
}
