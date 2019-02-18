package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"unsafe"

	"github.com/notti/nocgo"
)

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
	ret int `nocgo:"ret"`
}

type printCall struct {
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

	argsT := &testCall{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, -11, 12}
	argsP := &printCall{}

	var testvalue *int

	f1, err := l.Func("test_call", argsT)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(f1)

	f2, err := l.Func("print_value", argsP)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(argsT)
	f1.Call(unsafe.Pointer(argsT))
	fmt.Println(argsT)

	f2.Call(unsafe.Pointer(argsP))
	err = l.Value("value", &testvalue)
	if err != nil {
		log.Fatal(err)
	}
	*testvalue = 100
	f2.Call(unsafe.Pointer(argsP))
}
