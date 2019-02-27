/*
Package fakecgo fakes necessary functionality to support calling C-code (TLS initialization for main thread and subsequent threads).

Usage

Just import this library with
	import _ "github.com/notti/nocgo/fakecgo"
and enjoy the side effects (e.g., you can't use cgo any longer) :)
*/
package fakecgo

import (
	"syscall"
	"unsafe"
)

// WARNING: please read this before changing/improving anything
// This here might look like (ugly) go - but it is actually somehow C-code written in go (basically stuff in runtime/cgo/)
// Yeah this somehow works, but needs the trampolines from trampoline_*.s to fix calling conventions cdecl <-> go
//
// Beware that strings must end with a 0 byte to not confuse C-code we call
//
// Write barriers (we will be called while go is in a state where this is not possible) and stack split (we will be on systemstack anyway) are NOT allowed in here
// -> e.g. use memmove for copying to pointers
// go:nowritebarrierrec is only allowed inside runtime - so this has to be checked manually :(

// pthread_create will call this with ts as argument in a new thread -> this fixes up arguments to go (in assembly) and calls threadentry
func threadentry_trampoline()

// here we will store a pointer to the provided setg func
var setg_func uintptr

// x_cgo_init(G *g, void (*setg)(void*)) (runtime/cgo/gcc_linux_amd64.c)
// This get's called during startup, adjusts stacklo, and provides a pointer to setg_gcc for us
// Additionally, if we set _cgo_init to non-null, go won't do it's own TLS setup
//go:nosplit
func x_cgo_init(g *g, setg uintptr) {
	var size size_t
	var attr pthread_attr

	// we need an extra variable here - otherwise go generates "safe" code, which is not allowed here
	stackSp := uintptr(unsafe.Pointer(&size))

	setg_func = setg

	pthread_attr_init(&attr)
	pthread_attr_getstacksize(&attr, &size)
	g.stack.lo = stackSp - uintptr(size) + 4096
	pthread_attr_destroy(&attr)
}

// _cgo_thread_start is split into three parts in cgo since only one part is system dependent (keep it here for easier handling)

// _cgo_thread_start(ThreadStart *arg) (runtime/cgo/gcc_util.c)
// This get's called instead of the go code for creating new threads
// -> pthread_* stuff is used, so threads are setup correctly for C
// If this is missing, TLS is only setup correctly on thread 1!
//
//go:nosplit
func x_cgo_thread_start(arg *threadstart) {
	ts := malloc(unsafe.Sizeof(threadstart{}))

	if ts == 0 {
		dprintf(2, "couldn't allocate memory for threadstart\n\000")
		abort()
	}

	// *ts = *arg would cause a write barrier, which is not allowed
	memmove(unsafe.Pointer(ts), unsafe.Pointer(arg), unsafe.Sizeof(threadstart{}))

	_cgo_sys_start_thread((*threadstart)(unsafe.Pointer(ts)))
}

//go:nosplit
func _cgo_sys_start_thread(ts *threadstart) {
	var attr pthread_attr
	var ign, oset sigset_t
	var p pthread_t
	var size size_t

	sigfillset(&ign)
	pthread_sigmask(SIG_SETMASK, &ign, &oset)

	pthread_attr_init(&attr)
	pthread_attr_getstacksize(&attr, &size)
	(*g)(ts.g).stack.hi = uintptr(size)
	err := _cgo_try_pthread_create(&p, &attr, unsafe.Pointer(funcPC(threadentry_trampoline)), unsafe.Pointer(ts))

	pthread_sigmask(SIG_SETMASK, &oset, nil)

	if err != 0 {
		dprintf(2, "pthread_create failed: %s\n\000", strerror(err))
		abort()
	}
}

//go:nosplit
func _cgo_try_pthread_create(thread *pthread_t, attr *pthread_attr, start, arg unsafe.Pointer) int32 {
	for tries := 0; tries < 20; tries++ {
		err := pthread_create(thread, attr, start, arg)
		if err == 0 {
			pthread_detach(*thread)
			return 0
		}
		if err != int32(syscall.EAGAIN) {
			return err
		}
		ts := timespec{tv_sec: 0, tv_nsec: (tries + 1) * 1000 * 1000}
		nanosleep(&ts, nil)
	}
	return int32(syscall.EAGAIN)
}

func setg_trampoline(uintptr, unsafe.Pointer)

//go:nosplit
func threadentry(v unsafe.Pointer) uintptr {
	ts := *(*threadstart)(v)

	free(v)

	setg_trampoline(setg_func, ts.g)

	// faking funcs in go is a bit a... involved - but the following works :)
	fn := uintptr(unsafe.Pointer(&ts.fn))
	(*(*func())(unsafe.Pointer(&fn)))()

	return 0
}

// The following functions are required by the runtime - otherwise it complains via panic that they are missing

// do nothing - we don't support being a library for now
// _cgo_notify_runtime_init_done (runtime/cgo/gcc_libinit.c)
//go:nosplit
func x_cgo_notify_runtime_init_done() {}

// _cgo_setenv(char **arg) (runtime/cgo/gcc_setenv.c)
//go:nosplit
func x_cgo_setenv(arg [2]*byte) {
	setenv(arg[0], arg[1], 1)
}

// _cgo_unsetenv(char *arg) (runtime/cgo/gcc_setenv.c)
//go:nosplit
func x_cgo_unsetenv(arg *byte) {
	unsetenv(arg)
}
