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
