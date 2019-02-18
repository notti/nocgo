package nocgo 

import (
	"reflect"
	"unsafe"
)

var sliceOffset = reflect.TypeOf(reflect.SliceHeader{}).Field(0).Offset

// Call calls the given spec with the given arguments
func (spec Spec) Call(args unsafe.Pointer) {
	spec.base = uintptr(args)

	entersyscall()
	asmcgocall(unsafe.Pointer(asmcallptr), uintptr(unsafe.Pointer(&spec)))
	exitsyscall()

	if _Cgo_always_false {
		_Cgo_use(args)
		_Cgo_use(spec)
	}
}

//go:linkname _Cgo_always_false runtime.cgoAlwaysFalse
var _Cgo_always_false bool

//go:linkname _Cgo_use runtime.cgoUse
func _Cgo_use(interface{})

//go:linkname asmcgocall runtime.asmcgocall
func asmcgocall(unsafe.Pointer, uintptr) int32

//go:linkname entersyscall runtime.entersyscall
func entersyscall()

//go:linkname exitsyscall runtime.exitsyscall
func exitsyscall()

func asmcall()

//go:linkname x_cgo_init x_cgo_init
func x_cgo_init()

// force _cgo_init into the .data segment (instead of .bss), so our "linker" can overwrite its contents
//go:linkname _cgo_init _cgo_init
var _cgo_init = uintptr(10)

func init() {
	if _Cgo_always_false {
		x_cgo_init() // prevent x_cgo_init from being optimized out
	}
}

//go:linkname funcPC runtime.funcPC
func funcPC(f interface{}) uintptr

var asmcallptr = funcPC(asmcall)
