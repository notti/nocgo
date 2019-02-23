package ffi

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

// Call calls the given spec with the given arguments
func (spec Spec) Call(args unsafe.Pointer) {
	spec.base = uintptr(args)

	cgocall(unsafe.Pointer(asmcallptr), unsafe.Pointer(&spec))

	if _Cgo_always_false {
		_Cgo_use(spec)
	}
}

//go:linkname _Cgo_always_false runtime.cgoAlwaysFalse
var _Cgo_always_false bool

//go:linkname _Cgo_use runtime.cgoUse
func _Cgo_use(interface{})

// force _cgo_*, iscgo into the .data segment (instead of .bss), so our "linker" can overwrite its contents

//go:linkname x_cgo_init x_cgo_init
func x_cgo_init()

//go:linkname _cgo_init _cgo_init
var _cgo_init = uintptr(10)

//go:linkname x_cgo_thread_start x_cgo_thread_start
func x_cgo_thread_start()

//go:linkname _cgo_thread_start _cgo_thread_start
var _cgo_thread_start = uintptr(10)

//go:linkname x_cgo_notify_runtime_init_done x_cgo_notify_runtime_init_done
func x_cgo_notify_runtime_init_done()

//go:linkname _cgo_notify_runtime_init_done _cgo_notify_runtime_init_done
var _cgo_notify_runtime_init_done = uintptr(10)

//go:linkname x_cgo_setenv x_cgo_setenv
func x_cgo_setenv()

//go:linkname _cgo_setenv runtime._cgo_setenv
var _cgo_setenv = uintptr(10)

//go:linkname x_cgo_unsetenv x_cgo_unsetenv
func x_cgo_unsetenv()

//go:linkname _cgo_unsetenv runtime._cgo_unsetenv
var _cgo_unsetenv = uintptr(10)

//go:linkname x_cgo_callers x_cgo_callers
func x_cgo_callers()

//go:linkname _cgo_callers _cgo_callers
var _cgo_callers = uintptr(10)

//go:linkname iscgo runtime.iscgo
var iscgo = 1

func init() {
	if _Cgo_always_false {
		// prevent x_cgo_* from being optimized out
		x_cgo_init()
		x_cgo_thread_start()
		x_cgo_notify_runtime_init_done()
		x_cgo_setenv()
		x_cgo_unsetenv()
		x_cgo_callers()
	}
}
