package fakecgo

import "unsafe"

//go:linkname memmove runtime.memmove
func memmove(to, from unsafe.Pointer, n uintptr)

//go:linkname funcPC runtime.funcPC
func funcPC(f interface{}) uintptr

// the following struct is cgothreadstart from runtime
type threadstart struct {
	g   unsafe.Pointer //should be guintptr
	tls *uint64
	fn  unsafe.Pointer
}

// just enough from the runtime to manipulate g->stack->lo/hi
type stack struct {
	lo uintptr
	hi uintptr
}

type g struct {
	stack stack
}

// We actually don't need the full thing - but this is the same as in runtime and makes possible integration simpler
type libcall struct {
	fn   uintptr
	n    uintptr
	args uintptr
	r1   uintptr
	r2   uintptr
	err  uintptr
}

type libcFunc uintptr
