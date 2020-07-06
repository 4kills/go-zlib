package native

/*
#cgo CFLAGS: -I${SRCDIR}/zlib/
#cgo windows LDFLAGS: ${SRCDIR}/libs/winlibz.a
#cgo linux LDFLAGS: ${SRCDIR}/libs/linuxlibz.a
#cgo darwin LDFLAGS: ${SRCDIR}/libs/darwinlibz.a
*/
import "C"
