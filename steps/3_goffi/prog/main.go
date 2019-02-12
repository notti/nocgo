package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/notti/go-dynamic/steps/3_goffi/ffi"
)

var puts__dynload uintptr
var strcat__dynload uintptr

type putsString struct {
	s   []byte
	num int `ffi:"ret"`
}

func main() {
	fmt.Println(os.Args) // check if startup works

	str := "hello world"
	b := append([]byte(str), 0)

	args := &putsString{s: b}
	spec := ffi.MakeSpec(puts__dynload, args)

	fmt.Println(args, spec)

	spec.Call(unsafe.Pointer(args))

	fmt.Println(args)
}
