package main

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

var puts__dynload uintptr
var strcat__dynload uintptr

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

var dings uintptr

func main() {
	str := "hello world"
	b := append([]byte(str), 0)

	fmt.Println(call3(puts__dynload, uintptr(unsafe.Pointer(&b[0])), 0, 0))

	runtime.KeepAlive(b)

	teststr := []byte("hello\000world\000")
	teststr2 := []byte("C!\000")
	fmt.Println(call3(strcat__dynload, uintptr(unsafe.Pointer(&teststr[0])), uintptr(unsafe.Pointer(&teststr2[0])), 0))
	runtime.KeepAlive(teststr2)

	fmt.Println(call3(puts__dynload, uintptr(unsafe.Pointer(&teststr[0])), 0, 0))
	runtime.KeepAlive(teststr)

	fmt.Println(teststr)
}
