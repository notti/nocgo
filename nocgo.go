package nocgo 

import (
	"reflect"
	"unsafe"
)

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

//go:linkname _cgo_init _cgo_init
//go:linkname x_cgo_init x_cgo_init
var x_cgo_init byte
var _cgo_init = &x_cgo_init

//go:linkname x_cgo_thread_start x_cgo_thread_start
//go:linkname _cgo_thread_start _cgo_thread_start
var x_cgo_thread_start byte
var _cgo_thread_start = &x_cgo_thread_start

//go:linkname x_cgo_notify_runtime_init_done x_cgo_notify_runtime_init_done
//go:linkname _cgo_notify_runtime_init_done _cgo_notify_runtime_init_done
var x_cgo_notify_runtime_init_done byte
var _cgo_notify_runtime_init_done = &x_cgo_notify_runtime_init_done

//go:linkname x_cgo_setenv x_cgo_setenv
//go:linkname _cgo_setenv runtime._cgo_setenv
var x_cgo_setenv byte
var _cgo_setenv = &x_cgo_setenv

//go:linkname x_cgo_unsetenv x_cgo_unsetenv
//go:linkname _cgo_unsetenv runtime._cgo_unsetenv
var x_cgo_unsetenv byte
var _cgo_unsetenv = &x_cgo_unsetenv

//go:linkname x_cgo_callers x_cgo_callers
//go:linkname _cgo_callers _cgo_callers
var x_cgo_callers byte
var _cgo_callers = &x_cgo_callers

//go:linkname _iscgo runtime.iscgo
var _iscgo = true
