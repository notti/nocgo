package nocgo

import _ "unsafe" //needed for go:linkname

// The following "fakes" all the necessary stuff for pretending we're using cgo, without actually doing that
// -> iscgo will be set to true and all functions that are then required by the runtime implemented
// This is necessary to get TLS working in the mainthread (cgo_init) and in all other threads (cgo_thread_start).
// If we leave this out, libc can't use TLS since go runtime overwrites it (printf with %f already needs that)

// The actual functions are implemented in assembly (cgo_*.s)

//go:linkname _cgo_init _cgo_init
//go:linkname x_cgo_init x_cgo_init
var x_cgo_init byte
var _cgo_init = &x_cgo_init

//go:linkname x_cgo_thread_start x_cgo_thread_start
//go:linkname _cgo_thread_start _cgo_thread_start
var x_cgo_thread_start byte
var _cgo_thread_start = &x_cgo_thread_start

//go:linkname x_cgo_notify_runtime_init_done x_cgo_notify_runtime_init_done
//go:linkname _cgo_notify_runtime_init_done _cgo_notify_runtime_init_done
var x_cgo_notify_runtime_init_done byte
var _cgo_notify_runtime_init_done = &x_cgo_notify_runtime_init_done

//go:linkname x_cgo_setenv x_cgo_setenv
//go:linkname _cgo_setenv runtime._cgo_setenv
var x_cgo_setenv byte
var _cgo_setenv = &x_cgo_setenv

//go:linkname x_cgo_unsetenv x_cgo_unsetenv
//go:linkname _cgo_unsetenv runtime._cgo_unsetenv
var x_cgo_unsetenv byte
var _cgo_unsetenv = &x_cgo_unsetenv

//go:linkname x_cgo_callers x_cgo_callers
//go:linkname _cgo_callers _cgo_callers
var x_cgo_callers byte
var _cgo_callers = &x_cgo_callers

//go:linkname _iscgo runtime.iscgo
var _iscgo = true

// Now all the symbols we need to import from various libraries to implement the above functions:
// (__dynload is an artifact of relink.go)

//go:cgo_import_dynamic libc_pthread_attr_init pthread_attr_init "libpthread.so"
//go:linkname libc_pthread_attr_init libc_pthread_attr_init
//go:linkname pthread_attr_init__dynload pthread_attr_init__dynload
var libc_pthread_attr_init byte
var pthread_attr_init__dynload = &libc_pthread_attr_init

//go:cgo_import_dynamic libc_pthread_attr_getstacksize pthread_attr_getstacksize "libpthread.so"
//go:linkname libc_pthread_attr_getstacksize libc_pthread_attr_getstacksize
//go:linkname pthread_attr_getstacksize__dynload pthread_attr_getstacksize__dynload
var libc_pthread_attr_getstacksize byte
var pthread_attr_getstacksize__dynload = &libc_pthread_attr_getstacksize

//go:cgo_import_dynamic libc_pthread_attr_destroy pthread_attr_destroy "libpthread.so"
//go:linkname libc_pthread_attr_destroy libc_pthread_attr_destroy
//go:linkname pthread_attr_destroy__dynload pthread_attr_destroy__dynload
var libc_pthread_attr_destroy byte
var pthread_attr_destroy__dynload = &libc_pthread_attr_destroy

//go:cgo_import_dynamic libc_pthread_sigmask pthread_sigmask "libpthread.so"
//go:linkname libc_pthread_sigmask libc_pthread_sigmask
//go:linkname pthread_sigmask__dynload pthread_sigmask__dynload
var libc_pthread_sigmask byte
var pthread_sigmask__dynload = &libc_pthread_sigmask

//go:cgo_import_dynamic libc_pthread_create pthread_create "libpthread.so"
//go:linkname libc_pthread_create libc_pthread_create
//go:linkname pthread_create__dynload pthread_create__dynload
var libc_pthread_create byte
var pthread_create__dynload = &libc_pthread_create

//go:cgo_import_dynamic libc_pthread_detach pthread_detach "libpthread.so"
//go:linkname libc_pthread_detach libc_pthread_detach
//go:linkname pthread_detach__dynload pthread_detach__dynload
var libc_pthread_detach byte
var pthread_detach__dynload = &libc_pthread_detach

//go:cgo_import_dynamic libc_setenv setenv "libc.so.6"
//go:linkname libc_setenv libc_setenv
//go:linkname setenv__dynload setenv__dynload
var libc_setenv byte
var setenv__dynload = &libc_setenv

//go:cgo_import_dynamic libc_unsetenv unsetenv "libc.so.6"
//go:linkname libc_unsetenv libc_unsetenv
//go:linkname unsetenv__dynload unsetenv__dynload
var libc_unsetenv byte
var unsetenv__dynload = &libc_unsetenv

//go:cgo_import_dynamic libc_malloc malloc "libc.so.6"
//go:linkname libc_malloc libc_malloc
//go:linkname malloc__dynload malloc__dynload
var libc_malloc byte
var malloc__dynload = &libc_malloc

//go:cgo_import_dynamic libc_free free "libc.so.6"
//go:linkname libc_free libc_free
//go:linkname free__dynload free__dynload
var libc_free byte
var free__dynload = &libc_free

//go:cgo_import_dynamic libc_nanosleep nanosleep "libc.so.6"
//go:linkname libc_nanosleep libc_nanosleep
//go:linkname nanosleep__dynload nanosleep__dynload
var libc_nanosleep byte
var nanosleep__dynload = &libc_nanosleep

//go:cgo_import_dynamic libc_sigfillset sigfillset "libc.so.6"
//go:linkname libc_sigfillset libc_sigfillset
//go:linkname sigfillset__dynload sigfillset__dynload
var libc_sigfillset byte
var sigfillset__dynload = &libc_sigfillset

// on amd64 we don't need the following lines - on 386 we do...
// anyway - with those lines the output is better (but doesn't matter) - without it on amd64 we get multiple DT_NEEDED with "libc.so.6" etc

//go:cgo_import_dynamic _ _ "libpthread.so"
//go:cgo_import_dynamic _ _ "libc.so.6"
