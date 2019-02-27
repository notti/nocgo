package fakecgo

import _ "unsafe" // for go:linkname

// on BSDs we need the following (this is from runtime/cgo/freebsd.go)... which is not allowed outside cgo/stdlib :( - no way around that
// Right now we can fix that with -gcflags=github.com/notti/nocgo/fakecgo=-std during build - but that doesn't seem to work with go test

//go:linkname _environ environ
//go:linkname _progname __progname

//go:cgo_export_dynamic environ
//go:cgo_export_dynamic __progname

var _environ uintptr
var _progname uintptr

// we have to use the 3 argument format here :( - 2 argument format is only allowed from inside cgo

//go:cgo_import_dynamic libc_pthread_attr_init_x pthread_attr_init "libpthread.so"
//go:cgo_import_dynamic libc_pthread_attr_getstacksize_x pthread_attr_getstacksize "libpthread.so"
//go:cgo_import_dynamic libc_pthread_attr_destroy_x pthread_attr_destroy "libpthread.so"
//go:cgo_import_dynamic libc_pthread_sigmask_x pthread_sigmask "libpthread.so"
//go:cgo_import_dynamic libc_pthread_create_x pthread_create "libpthread.so"
//go:cgo_import_dynamic libc_pthread_detach_x pthread_detach "libpthread.so"
//go:cgo_import_dynamic libc_setenv_x setenv "libc.so.7"
//go:cgo_import_dynamic libc_unsetenv_x unsetenv "libc.so.7"
//go:cgo_import_dynamic libc_malloc_x malloc "libc.so.7"
//go:cgo_import_dynamic libc_free_x free "libc.so.7"
//go:cgo_import_dynamic libc_nanosleep_x nanosleep "libc.so.7"
//go:cgo_import_dynamic libc_sigfillset_x sigfillset "libc.so.7"
//go:cgo_import_dynamic libc_abort_x abort "libc.so.7"
//go:cgo_import_dynamic libc_dprintf_x dprintf "libc.so.7"
//go:cgo_import_dynamic libc_strerror_x strerror "libc.so.7"
//go:cgo_import_dynamic _ _ "libpthread.so"
//go:cgo_import_dynamic _ _ "libc.so.7"
