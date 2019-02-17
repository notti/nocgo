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
	i1  uint16
	i2  int
	f1  float32
	f2  float64
	i3  int
	i4  int
	i5  int
	i6  int
	i7  int
	i8  int
	i9  int16
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

	argsT := &testCall{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, -11, 12}
	specT := ffi.MakeSpec(test_call__dynload, argsT)

	fmt.Println(argsT, specT)
	specT.Call(unsafe.Pointer(argsT))
	fmt.Println(argsT)
}
