package fakecgo

import (
	"reflect"
	"unsafe"
)

// Wrapper functions to provide libc-functions for cgo.go

// C-types (linux glibc, libc freebsd) we don't need to see what's inside pthread_attr and sigset_t -> byte arrays:
type pthread_attr [56]byte // this is 36 on 386 glibc, 16 on amd64 freebsd - but too big doesn't hurt
type sigset_t [128]byte    // this is 16 on amd64 freebsd - but too big doesn't hurt
type size_t int
type pthread_t int
type timespec struct {
	tv_sec  int
	tv_nsec int
}

// We could take timespec from syscall - but there it uses int32 and int64 for 32 bit and 64 bit arch, which complicates stuff for us

// for pthread_sigmask:

type sighow int32

const (
	SIG_BLOCK   sighow = 0
	SIG_UNBLOCK sighow = 1
	SIG_SETMASK sighow = 2
)

// Every wrapper here MUST NOT split stack or have write barriers! (see cgoGlue.go)

//go:nosplit
func pthread_attr_init(attr *pthread_attr) int32 {
	return int32(libcCall1(libc_pthread_attr_init, uintptr(unsafe.Pointer(attr))))
}

//go:nosplit
func pthread_attr_getstacksize(attr *pthread_attr, stacksize *size_t) int32 {
	return int32(libcCall2(libc_pthread_attr_getstacksize, uintptr(unsafe.Pointer(attr)), uintptr(unsafe.Pointer(stacksize))))
}

//go:nosplit
func pthread_attr_destroy(attr *pthread_attr) int32 {
	return int32(libcCall1(libc_pthread_attr_destroy, uintptr(unsafe.Pointer(attr))))
}

//go:nosplit
func pthread_sigmask(how sighow, set, oldset *sigset_t) int32 {
	return int32(libcCall3(libc_pthread_sigmask, uintptr(how), uintptr(unsafe.Pointer(set)), uintptr(unsafe.Pointer(oldset))))
}

//go:nosplit
func pthread_create(thread *pthread_t, attr *pthread_attr, start, arg unsafe.Pointer) int32 {
	return int32(libcCall4(libc_pthread_create, uintptr(unsafe.Pointer(thread)), uintptr(unsafe.Pointer(attr)), uintptr(start), uintptr(arg)))
}

//go:nosplit
func pthread_detach(thread pthread_t) int32 {
	return int32(libcCall1(libc_pthread_detach, uintptr(thread)))
}

//go:nosplit
func setenv(name, value *byte, overwrite int32) int32 {
	return int32(libcCall3(libc_setenv, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(value)), uintptr(overwrite)))
}

//go:nosplit
func unsetenv(name *byte) int32 {
	return int32(libcCall1(libc_unsetenv, uintptr(unsafe.Pointer(name))))
}

//go:nosplit
func malloc(size uintptr) uintptr {
	return libcCall1(libc_malloc, uintptr(size))
}

//go:nosplit
func free(ptr unsafe.Pointer) {
	libcCall1(libc_free, uintptr(ptr))
}

//go:nosplit
func nanosleep(rgtp, rmtp *timespec) int32 {
	return int32(libcCall2(libc_nanosleep, uintptr(unsafe.Pointer(rgtp)), uintptr(unsafe.Pointer(rmtp))))
}

//go:nosplit
func sigfillset(set *sigset_t) int32 {
	return int32(libcCall1(libc_sigfillset, uintptr(unsafe.Pointer(set))))
}

//go:nosplit
func abort() {
	libcCall0(libc_abort)
}

//go:nosplit
func dprintf(fd uintptr, fmt string, arg ...uintptr) int32 {
	var args [4]uintptr
	copy(args[:], arg)
	return int32(libcCall6(libc_dprintf, fd, (*reflect.StringHeader)(unsafe.Pointer(&fmt)).Data, args[0], args[1], args[2], args[3]))
}

//go:nosplit
func strerror(err int32) uintptr {
	return libcCall1(libc_strerror, uintptr(err))
}
