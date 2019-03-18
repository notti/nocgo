package nocgo

import (
	"errors"
	"reflect"
	"unsafe"
)

func mustSpec(fn *byte, fun interface{}) {
	err := makeSpec(uintptr(unsafe.Pointer(fn)), fun)
	if err != nil {
		panic(err)
	}
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

var dlopen func(filename []byte, flags int32) uintptr
var dlclose func(handle uintptr) int32
var dlsym func(handle uintptr, symbol []byte) uintptr
var dlerror func() uintptr

func init() {
	mustSpec(libc_dlopen, &dlopen)
	mustSpec(libc_dlclose, &dlclose)
	mustSpec(libc_dlsym, &dlsym)
	mustSpec(libc_dlerror, &dlerror)
}

func getLastError() error {
	err := dlerror()
	if err == 0 {
		return errors.New("Unknown dl error")
	}
	return errors.New(MakeGoStringFromPointer(err))
}

// Library holds the handle to a loaded library
type Library uintptr

// Open opens the given dynamic library and returns a handle for loading symbols and functions.
func Open(library string) (Library, error) {
	handle := dlopen(MakeCString(library), 2 /* RTLD_NOW */)
	if handle != 0 {
		return Library(handle), nil
	}
	return 0, getLastError()
}

// Close closes the library. This might also release all resources. Any Func and Value calls on the Library after this point can give unexpected results.
func (l Library) Close() error {
	ret := dlclose(uintptr(l))
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
	addr := dlsym(uintptr(l), MakeCString(name))
	if addr == 0 {
		return getLastError()
	}
	return makeSpec(addr, fun)
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

	addr := dlsym(uintptr(l), MakeCString(name))
	if addr == 0 {
		return getLastError()
	}

	*(*uintptr)(unsafe.Pointer(v.UnsafeAddr())) = addr

	return nil
}
