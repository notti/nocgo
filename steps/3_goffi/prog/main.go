package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/notti/go-dynamic/steps/3_goffi/ffi"
)

var puts__dynload uintptr
var test_call__dynload uintptr

type putsString struct {
	s   []byte
	num int `ffi:"ret"`
}

type testCall struct {
	a   uint16
	b   int
	c   float32
	d   float64
	ret int `ffi:"ret"`
}

func main() {
	fmt.Println(os.Args) // check if startup works

	str := "hello world"
	b := append([]byte(str), 0)

	argsP := &putsString{s: b}
	specP := ffi.MakeSpec(puts__dynload, argsP)

	fmt.Println(argsP, specP)
	specP.Call(unsafe.Pointer(argsP))
	fmt.Println(argsP)

	argsT := &testCall{1, 2, 3, 4, 5}
	specT := ffi.MakeSpec(test_call__dynload, argsT)

	fmt.Println(argsT, specT)
	specT.Call(unsafe.Pointer(argsT))
	fmt.Println(argsT)
}
