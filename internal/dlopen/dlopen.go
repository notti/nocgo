// +build !cgo

package dlopen

import (
	"unsafe"

	_ "github.com/notti/nocgo/internal/fakecgo" // get everything we need to fake cgo behaviour if we're not using cgo
	"github.com/notti/nocgo/internal/ffi"
)

func mustSpec(fn *byte, fun interface{}) {
	err := ffi.MakeSpec(uintptr(unsafe.Pointer(fn)), fun)
	if err != nil {
		panic(err)
	}
}

// on 386 we need to do the dance of cgo_import_dynamic followed by two linknames,
// definining a variable that gets the dynamic symbol, and then derefercing it.
// Othwerwise we get an unknown relocation type error during linking

//go:linkname libc_dlopen_x libc_dlopen_x
var libc_dlopen_x byte
var libc_dlopen = &libc_dlopen_x

//go:linkname libc_dlclose_x libc_dlclose_x
var libc_dlclose_x byte
var libc_dlclose = &libc_dlclose_x

//go:linkname libc_dlsym_x libc_dlsym_x
var libc_dlsym_x byte
var libc_dlsym = &libc_dlsym_x

//go:linkname libc_dlerror_x libc_dlerror_x
var libc_dlerror_x byte
var libc_dlerror = &libc_dlerror_x

var DLOpen func(filename []byte, flags int32) uintptr
var DLClose func(handle uintptr) int32
var DLSym func(handle uintptr, symbol []byte) uintptr
var DLError func() uintptr

func init() {
	mustSpec(libc_dlopen, &DLOpen)
	mustSpec(libc_dlclose, &DLClose)
	mustSpec(libc_dlsym, &DLSym)
	mustSpec(libc_dlerror, &DLError)
}
