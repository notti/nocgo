package nocgo

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/notti/nocgo/internal/dlopen"
	"github.com/notti/nocgo/internal/ffi"
)

func getLastError() error {
	err := dlopen.DLError()
	if err == 0 {
		return errors.New("Unknown dl error")
	}
	return errors.New(MakeGoStringFromPointer(err))
}

// Library holds the handle to a loaded library
type Library uintptr

// Open opens the given dynamic library and returns a handle for loading symbols and functions.
func Open(library string) (Library, error) {
	handle := dlopen.DLOpen(MakeCString(library), 2 /* RTLD_NOW */)
	if handle != 0 {
		return Library(handle), nil
	}
	return 0, getLastError()
}

// Close closes the library. This might also release all resources. Any Func and Value calls on the Library after this point can give unexpected results.
func (l Library) Close() error {
	ret := dlopen.DLClose(uintptr(l))
	if ret == 0 {
		return nil
	}
	return getLastError()
}

// Func returns a callable spec for the given symbol name and argument specification.
//
// WARNING! This does not and cannot check if the size of the given type is correct!
//
// Example:
//	var puts func(s []byte) int32
//	if err := lib.Func("puts", &puts); err != nil {
//		//handle error; err will contain an error message from dlerror, or if something went wrong with building the spec
//	}
//	num := puts(nocgo.MakeCString("hello world!\n"))
//	fmt.Printf("Successfully printed %d characters from C!\n", num)
//
// See package documentation for an explanation of C-types
func (l Library) Func(name string, fun interface{}) error {
	addr := dlopen.DLSym(uintptr(l), MakeCString(name))
	if addr == 0 {
		return getLastError()
	}
	return ffi.MakeSpec(addr, fun)
}

// Value sets the given value (which must be pointer to pointer to the correct type) to the global symbol given by name.
//
// WARNING! This does not and cannot check if the size of the given type is correct! This might be possibly dangerous.
// See above for an explanation of C-types.
//
// Example:
// 	var value *int32
//	if err := lib.Value("some_value", &value); err != nil {
//		//handle error; error will contain an error message from dlerror
//	}
//
//	// *value now is the contents of the global symbol in the library
//	fmt.Printf(*value)
func (l Library) Value(name string, value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return errors.New("value must be a pointer to a pointer")
	}
	v = v.Elem()
	if v.Kind() != reflect.Ptr {
		return errors.New("value must be pointer to a pointer")
	}

	addr := dlopen.DLSym(uintptr(l), MakeCString(name))
	if addr == 0 {
		return getLastError()
	}

	*(*uintptr)(unsafe.Pointer(v.UnsafeAddr())) = addr

	return nil
}
