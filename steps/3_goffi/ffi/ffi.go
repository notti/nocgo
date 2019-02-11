package ffi

import (
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

type buildType int

const (
	typeCstring buildType = iota
)

type builder struct {
	index int
	kind  buildType
}

type argumentType uintptr

const (
	typePointer argumentType = 0
	typeInteger argumentType = 1
)

type argument struct {
	offset uintptr
	kind   argumentType
}

// Spec is the callspec needed to do the actuall call
type Spec struct {
	fn      uintptr
	base    uintptr
	regargs [6]uintptr
	rax     uintptr
	ret0    uintptr
	ret1    uintptr
}

var sliceOffset = reflect.TypeOf(reflect.SliceHeader{}).Field(0).Offset

func fieldToOffset(k reflect.StructField, t string) uintptr {
	switch k.Type.Kind() {
	case reflect.Slice:
		return k.Offset + sliceOffset // FIXME: is this correct?
	case reflect.Int:
		return k.Offset
		// TODO: implement types
	}
	panic("Unknown Type")
}

func MakeSpec(args interface{}) Spec {
	v := reflect.ValueOf(args)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	var spec Spec
	haveRet := false

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tags := strings.Split(f.Tag.Get("ffi"), ",")
		ret := false
		st := ""
		for _, tag := range tags {
			if tag == "ret" {
				if haveRet == true {
					panic("Only one return argument allowed")
				}
				ret = true
				haveRet = true
				continue
			}
			if strings.HasPrefix(tag, "type=") {
				st = tag[5:]
			}
		}
		if ret {
			spec.ret0 = fieldToOffset(f, st)
			// FIXME ret1!
			continue
		}
		spec.regargs[i] = fieldToOffset(f, st)
		// FIXME stackspill, float, xmm
	}
	// FIXME set rax
	return spec
}

/*
	fn
	base
	regarg0 rdi
	regarg1 rsi
	regarg2 rdx
	regarg3 rcx
	regarg4 r8
	regarg5 r9
	rax     rax
	ret0    rax
	ret1    rdx
*/

var asmcall3ptr = unsafe.Pointer(reflect.ValueOf(asmcall3).Pointer())

func Call(fn uintptr, spec Spec, args unsafe.Pointer) {

	spec.fn = fn
	spec.base = uintptr(args)

	entersyscall()
	asmcgocall(asmcall3ptr, uintptr(unsafe.Pointer(&spec)))
	exitsyscall()

	runtime.KeepAlive(args)
	runtime.KeepAlive(spec)
}

//go:linkname asmcgocall runtime.asmcgocall
func asmcgocall(unsafe.Pointer, uintptr) int32

//go:linkname entersyscall runtime.entersyscall
func entersyscall()

//go:linkname exitsyscall runtime.exitsyscall
func exitsyscall()

func asmcall3()

func call3(fn uintptr, arg0 uintptr, arg1 uintptr, arg2 uintptr) uintptr {
	p := unsafe.Pointer(reflect.ValueOf(asmcall3).Pointer())

	entersyscall()
	asmcgocall(p, uintptr(unsafe.Pointer(&fn)))
	exitsyscall()

	runtime.KeepAlive(p)
	runtime.KeepAlive(fn)
	return fn
}
