package nocgo 

import (
	"reflect"
	"unsafe"
)

// sliceOffset is needed by MakeSpec
var sliceOffset = reflect.TypeOf(reflect.SliceHeader{}).Field(0).Offset

//go:linkname funcPC runtime.funcPC
func funcPC(f interface{}) uintptr

func asmcall()

var asmcallptr = funcPC(asmcall)

//go:linkname cgocall runtime.cgocall
func cgocall(fn, arg unsafe.Pointer) int32

//go:linkname _Cgo_always_false runtime.cgoAlwaysFalse
var _Cgo_always_false bool

//go:linkname _Cgo_use runtime.cgoUse
func _Cgo_use(interface{})

// Call calls the given spec with the given arguments
func (spec Spec) Call(args unsafe.Pointer) {
	spec.base = uintptr(args)

	// noescape doesn't work outside runtime
	// we don't support callbacks now - so the stack can not move from under us -> we probably can get away with this (hopefully)
	// otherwise we get a heap allocation here and performance goes out the window
	//
	// this should be solve in a better way; e.g.
	// - use *Spec and return *Spec from MakeSpec (so Spec lives on the heap)
	// - pass args in g.m.libcall
	specNoescape := uintptr(unsafe.Pointer(&spec))

	cgocall(unsafe.Pointer(asmcallptr), unsafe.Pointer(specNoescape))

	if _Cgo_always_false {
		_Cgo_use(spec)
	}
}
