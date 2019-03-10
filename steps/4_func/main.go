package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

//go:linkname funcPC runtime.funcPC
func funcPC(f interface{}) uintptr

func callWrapper()

type funcStorage struct {
	code  uintptr
	value int
}

func fake(s *funcStorage, i *int) int {
	fmt.Println("in fake: ", s.value, *i)
	return *i + 10 + s.value
}

func emulate(x reflect.Value) {
	ftype := x.Elem().Type()
	fptr := x.Pointer()
	fmt.Println(fptr, ftype)
	toassign := new(funcStorage)
	toassign.code = funcPC(callWrapper)
	toassign.value = -5
	*(*unsafe.Pointer)(unsafe.Pointer(fptr)) = unsafe.Pointer(toassign)
	fmt.Println(toassign)
}

func main() {
	var test func(int) int

	emulate(reflect.ValueOf(&test))

	fmt.Println("testcall1: ", test(10))

	fmt.Println("testcall2: ", test(11))

	fmt.Println("testcall3: ", test(12))

}
