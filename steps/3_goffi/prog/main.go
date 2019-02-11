package main

import (
	"fmt"
	"unsafe"

	"github.com/notti/go-dynamic/steps/3_goffi/ffi"
)

var puts__dynload uintptr
var strcat__dynload uintptr

type putsString struct {
	s   []byte
	num int `ffi:"ret"`
}

/*
type putsBuffer struct {
	ffi.FunctionSpec
	s []byte
}*/

/*
type strcat struct {
	ffi.FunctionSpec
	dest, src []byte
}*/

var dings uintptr

func main() {
	str := "hello world"
	b := append([]byte(str), 0)

	args := &putsString{s: b}
	spec := ffi.MakeSpec(args)

	fmt.Println(args, spec)

	ffi.Call(puts__dynload, spec, unsafe.Pointer(args))

	fmt.Println(args)

	/*

		fmt.Println(call3(puts__dynload, uintptr(unsafe.Pointer(&b[0])), 0, 0))

		runtime.KeepAlive(b)

		teststr := []byte("hello\000world\000")
		teststr2 := []byte("C!\000")
		fmt.Println(call3(strcat__dynload, uintptr(unsafe.Pointer(&teststr[0])), uintptr(unsafe.Pointer(&teststr2[0])), 0))
		runtime.KeepAlive(teststr2)

		fmt.Println(call3(puts__dynload, uintptr(unsafe.Pointer(&teststr[0])), 0, 0))
		runtime.KeepAlive(teststr)

		fmt.Println(teststr)*/
}
