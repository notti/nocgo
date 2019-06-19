package fakecgo

import "unsafe"

// wrapper functions to provide go func definitions with arguments

func asmlibccall6(fn, n, args uintptr) uintptr

//go:nosplit
func libcCall0(fn *libcFunc) uintptr {
	return asmlibccall6(uintptr(unsafe.Pointer(fn)), 0, 0)
}

//go:nosplit
//go:noinline
func libcCall1(fn *libcFunc, a1 uintptr) uintptr {
	return asmlibccall6(uintptr(unsafe.Pointer(fn)), 1, uintptr(unsafe.Pointer(&a1)))
}

//go:nosplit
//go:noinline
func libcCall2(fn *libcFunc, a1 uintptr, a2 uintptr) uintptr {
	return asmlibccall6(uintptr(unsafe.Pointer(fn)), 2, uintptr(unsafe.Pointer(&a1)))
}

//go:nosplit
//go:noinline
func libcCall3(fn *libcFunc, a1 uintptr, a2 uintptr, a3 uintptr) uintptr {
	return asmlibccall6(uintptr(unsafe.Pointer(fn)), 3, uintptr(unsafe.Pointer(&a1)))
}

//go:nosplit
//go:noinline
func libcCall4(fn *libcFunc, a1 uintptr, a2 uintptr, a3 uintptr, a4 uintptr) uintptr {
	return asmlibccall6(uintptr(unsafe.Pointer(fn)), 4, uintptr(unsafe.Pointer(&a1)))
}

//go:nosplit
//go:noinline
func libcCall5(fn *libcFunc, a1 uintptr, a2 uintptr, a3 uintptr, a4 uintptr, a5 uintptr) uintptr {
	return asmlibccall6(uintptr(unsafe.Pointer(fn)), 5, uintptr(unsafe.Pointer(&a1)))
}

//go:nosplit
//go:noinline
func libcCall6(fn *libcFunc, a1 uintptr, a2 uintptr, a3 uintptr, a4 uintptr, a5 uintptr, a6 uintptr) uintptr {
	return asmlibccall6(uintptr(unsafe.Pointer(fn)), 6, uintptr(unsafe.Pointer(&a1)))
}
