//go:build !fake

package amtrpc

/*
#cgo LDFLAGS: -L'$ORIGIN' -lrpc -Wl,-rpath='$ORIGIN'
#include <stdlib.h>
#include "librpc.h"
*/
import "C" //nolint:typecheck
import "unsafe"

type AMTRPC struct{}

func (A AMTRPC) CheckAccess() int {
	return int(C.rpcCheckAccess())
}

func (A AMTRPC) Exec(command string) (string, int) {
	ccmd := C.CString(command)
	cresponse := C.CString("")
	cstatus := C.rpcExec(ccmd, &cresponse)

	C.free(unsafe.Pointer(ccmd))
	C.free(unsafe.Pointer(cresponse))

	return C.GoString(cresponse), int(cstatus)
}
