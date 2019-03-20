// Smoke test/example demonstrating accessing global variables and functions from a self compiled library.
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/notti/nocgo"
)

func cb(a int32) int32 {
	fmt.Println("In go: ", a)
	return a * 2
}

func main() {
	fmt.Println(os.Args) // check if startup works

	var lib string
	switch runtime.GOARCH {
	case "386":
		lib = "libcalltest32.so.1"
	case "amd64":
		lib = "libcalltest64.so.1"
	default:
		log.Fatalln("Unknown arch ", runtime.GOARCH)
	}

	l, err := nocgo.Open(lib)
	if err != nil {
		log.Fatal(err)
	}

	var testCall func(
		i1 uint16,
		i2 int32,
		f1 float32,
		f2 float64,
		i3 int32,
		i4 int32,
		i5 int32,
		i6 int32,
		i7 int32,
		i8 int32,
		i9 int16,
	) int32

	var printCall func()

	var testCB func(cb func(int32) int32) int32

	var testvalue *int

	if err := l.Func("test_call", &testCall); err != nil {
		log.Fatal(err)
	}

	if err := l.Func("print_value", &printCall); err != nil {
		log.Fatal(err)
	}

	if err := l.Func("test_cb", &testCB); err != nil {
		log.Fatal(err)
	}

	fmt.Println(testCall(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, -11))

	printCall()
	err = l.Value("value", &testvalue)
	if err != nil {
		log.Fatal(err)
	}
	*testvalue = 100
	printCall()

	fmt.Println("back in go: ", testCB(cb))
}
