package nocgo

import (
	"errors"
	"reflect"
	"unsafe"

	_ "github.com/notti/nocgo/fakecgo" // get everything we need to fake cgo behaviour
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

//go:cgo_import_dynamic libc_dlopen dlopen "libdl.so"
//go:linkname libc_dlopen libc_dlopen
//go:linkname dlopen__dynload dlopen__dynload
var libc_dlopen byte
var dlopen__dynload = &libc_dlopen

//go:cgo_import_dynamic libc_dlclose dlclose "libdl.so"
//go:linkname libc_dlclose libc_dlclose
//go:linkname dlclose__dynload dlclose__dynload
var libc_dlclose byte
var dlclose__dynload = &libc_dlclose

//go:cgo_import_dynamic libc_dlsym dlsym "libdl.so"
//go:linkname libc_dlsym libc_dlsym
//go:linkname dlsym__dynload dlsym__dynload
var libc_dlsym byte
var dlsym__dynload = &libc_dlsym

//go:cgo_import_dynamic libc_dlerror dlerror "libdl.so"
//go:linkname libc_dlerror libc_dlerror
//go:linkname dlerror__dynload dlerror__dynload
var libc_dlerror byte
var dlerror__dynload = &libc_dlerror

// on amd64 we don't need the following line - on 386 we do...
// anyway - with those lines the output is better (but doesn't matter) - without it on amd64 we get multiple DT_NEEDED with "libc.so.6" etc

//go:cgo_import_dynamic _ _ "libdl.so"

var dlopenSpec = mustSpec(dlopen__dynload, dlopen{})
var dlcloseSpec = mustSpec(dlclose__dynload, dlclose{})
var dlsymSpec = mustSpec(dlsym__dynload, dlsym{})
var dlerrorSpec = mustSpec(dlerror__dynload, dlerror{})

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
