package main

// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
import "C"

import "unsafe"

func main() {
	cs := []byte("Hello\000 world\000")
	C.puts((*C.char)(unsafe.Pointer(&cs[0])))
	C.fputs((*C.char)(unsafe.Pointer(&cs[0])), C.stdout)
	C.putc('a', C.stdout)
	C.strcat((*C.char)(unsafe.Pointer(&cs[0])), (*C.char)(unsafe.Pointer(&cs[0])))
	C.puts((*C.char)(unsafe.Pointer(&cs[0])))
}
