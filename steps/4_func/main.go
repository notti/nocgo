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
	code    uintptr
	argsize uintptr
	value   int
}

func fake(s *funcStorage, i *[2]int) int {
	fmt.Println("in fake: ", s.value, i[0], i[1])
	return i[0] + 100 + s.value
}

var alignSize = unsafe.Sizeof(uintptr(0))

func emulate(x reflect.Value) {
	ftype := x.Elem().Type()
	fptr := x.Pointer()
	fmt.Println(fptr, ftype)
	toassign := new(funcStorage)
	toassign.code = funcPC(callWrapper)
	for i := 0; i < ftype.NumIn(); i++ {
		size := ftype.In(i).Size()
		size = (size + alignSize - 1) &^ (alignSize - 1)
		toassign.argsize += size
	}
	toassign.value = -5
	*(*unsafe.Pointer)(unsafe.Pointer(fptr)) = unsafe.Pointer(toassign)
	fmt.Println(toassign)
}

func main() {
	var test func(int, byte, int) int
	var test1 func(int, int) int
	var test2 func(int, string) int

	emulate(reflect.ValueOf(&test))
	emulate(reflect.ValueOf(&test1))
	emulate(reflect.ValueOf(&test2))

	fmt.Println("test 1: ", test(10, 1, 2))
	fmt.Println("test 2: ", test(11, 2, 3))
	fmt.Println("test 3: ", test(12, 3, 4))

	fmt.Println("test 1: ", test1(10, 1))
	fmt.Println("test 2: ", test1(11, 2))
	fmt.Println("test 3: ", test1(12, 3))

	fmt.Println("test2 1: ", test2(10, "1, 2"))
	fmt.Println("test2 2: ", test2(11, "2, 3"))
	fmt.Println("test2 3: ", test2(12, "3, 4"))

}
