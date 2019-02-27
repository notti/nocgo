package nocgo

import (
	"errors"
	"reflect"
	"unsafe"
)

type dlopen struct {
	filename []byte
	flags    int32
	handle   uintptr `nocgo:"ret"`
}

type dlclose struct {
	handle uintptr
	ret    int32 `nocgo:"ret"`
}

type dlsym struct {
	handle uintptr
	symbol []byte
	addr   uintptr `nocgo:"ret"`
}

type dlerror struct {
	err uintptr `nocgo:"ret"`
}

func mustSpec(fn *byte, args interface{}) Spec {
	ret, err := makeSpec(uintptr(unsafe.Pointer(fn)), args)
	if err != nil {
		panic(err)
	}
	return ret
}

// on 386 we need to do the dance of cgo_import_dynamic followed by two linknames,
// definining a variable that gets the dynamic symbol, and then derefercing it.
// Othwerwise we get an unknown relocation type error during linking

//go:linkname libc_dlopen_x libc_dlopen_x
var libc_dlopen_x byte
var libc_dlopen = &libc_dlopen_x

//go:linkname libc_dlclose_x libc_dlclose_x
var libc_dlclose_x byte
var libc_dlclose = &libc_dlclose_x

//go:linkname libc_dlsym_x libc_dlsym_x
var libc_dlsym_x byte
var libc_dlsym = &libc_dlsym_x

//go:linkname libc_dlerror_x libc_dlerror_x
var libc_dlerror_x byte
var libc_dlerror = &libc_dlerror_x

var dlopenSpec = mustSpec(libc_dlopen, dlopen{})
var dlcloseSpec = mustSpec(libc_dlclose, dlclose{})
var dlsymSpec = mustSpec(libc_dlsym, dlsym{})
var dlerrorSpec = mustSpec(libc_dlerror, dlerror{})

func getLastError() error {
	args := dlerror{}
	dlerrorSpec.Call(unsafe.Pointer(&args))
	if args.err == 0 {
		return errors.New("Unknown dl error")
	}
	return errors.New(MakeGoStringFromPointer(args.err))
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
