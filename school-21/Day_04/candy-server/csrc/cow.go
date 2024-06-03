package csrc

/*
#include "cow.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// AskCow вызывает функцию C ask_cow
func AskCow(phrase string) string {
	cstr := C.CString(phrase)
	defer C.free(unsafe.Pointer(cstr))

	result := C.ask_cow(cstr)
	defer C.free(unsafe.Pointer(result))

	return C.GoString(result)
}
