package ffi

import (
	"errors"
	"reflect"
	"unsafe"
)

type dlopen struct {
	filename []byte
	flags    int32
	handle   uintptr `ffi:"ret"`
}

type dlclose struct {
	handle uintptr
	ret    int32 `ffi:"ret"`
}

type dlsym struct {
	handle uintptr
	symbol []byte
	addr   uintptr `ffi:"ret"`
}

type dlerror struct {
	err uintptr `ffi:"ret"`
}

func mustSpec(fn uintptr, args interface{}) Spec {
	ret, err := makeSpec(fn, args)
	if err != nil {
		panic(err)
	}
	return ret
}

var dlopenSpec = mustSpec(uintptr(unsafe.Pointer(dlopen__dynload)), dlopen{})
var dlcloseSpec = mustSpec(uintptr(unsafe.Pointer(dlclose__dynload)), dlclose{})
var dlsymSpec = mustSpec(uintptr(unsafe.Pointer(dlsym__dynload)), dlsym{})
var dlerrorSpec = mustSpec(uintptr(unsafe.Pointer(dlerror__dynload)), dlerror{})

func getLastError() error {
	args := dlerror{}
	dlerrorSpec.Call(unsafe.Pointer(&args))
	if args.err == 0 {
		return errors.New("Unknown dl error")
	}
	return errors.New(MakeGoString(args.err))
}

// Library holds loaded library
type Library uintptr

// Open opens the given dynamic library and returns a handle for loading symbols and functions
func Open(library string) (Library, error) {
	args := dlopen{
		filename: MakeCString(library),
		flags:    2, // RTLD_NOW
	}
	dlopenSpec.Call(unsafe.Pointer(&args))
	if args.handle != 0 {
		return Library(args.handle), nil
	}
	return 0, getLastError()
}

// Close closes the library
func (l Library) Close() error {
	args := dlclose{
		handle: uintptr(l),
	}
	dlcloseSpec.Call(unsafe.Pointer(&args))
	if args.ret == 0 {
		return nil
	}
	return getLastError()
}

// Func returns a callable spec for the given symbol name and argument specification
func (l Library) Func(name string, args interface{}) (Spec, error) {
	a := dlsym{
		handle: uintptr(l),
		symbol: MakeCString(name),
	}
	dlsymSpec.Call(unsafe.Pointer(&a))
	if a.addr == 0 {
		return Spec{}, getLastError()
	}
	return makeSpec(a.addr, args)
}

// Value sets the given value (which must be pointer to pointer to the correct type) to the symbol given by name
func (l Library) Value(name string, value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return errors.New("value must be a pointer to a pointer")
	}
	v = v.Elem()
	if v.Kind() != reflect.Ptr {
		return errors.New("value must be pointer to a pointer")
	}

	a := dlsym{
		handle: uintptr(l),
		symbol: MakeCString(name),
	}
	dlsymSpec.Call(unsafe.Pointer(&a))
	if a.addr == 0 {
		return getLastError()
	}

	*(*uintptr)(unsafe.Pointer(v.UnsafeAddr())) = a.addr

	return nil
}
