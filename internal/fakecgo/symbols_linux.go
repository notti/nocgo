package fakecgo

// we have to use the 3 argument format here :( - 2 argument format is only allowed from inside cgo

// pthread_attr_init will get us the wrong version on glibc - but this doesn't matter, since the memory we provide is zeroed - which will lead the correct result again

//go:cgo_import_dynamic libc_pthread_attr_init_x pthread_attr_init "libpthread.so.0"
//go:cgo_import_dynamic libc_pthread_attr_getstacksize_x pthread_attr_getstacksize "libpthread.so.0"
//go:cgo_import_dynamic libc_pthread_attr_destroy_x pthread_attr_destroy "libpthread.so.0"
//go:cgo_import_dynamic libc_pthread_sigmask_x pthread_sigmask "libpthread.so.0"
//go:cgo_import_dynamic libc_pthread_create_x pthread_create "libpthread.so.0"
//go:cgo_import_dynamic libc_pthread_detach_x pthread_detach "libpthread.so.0"
//go:cgo_import_dynamic libc_setenv_x setenv "libc.so.6"
//go:cgo_import_dynamic libc_unsetenv_x unsetenv "libc.so.6"
//go:cgo_import_dynamic libc_malloc_x malloc "libc.so.6"
//go:cgo_import_dynamic libc_free_x free "libc.so.6"
//go:cgo_import_dynamic libc_nanosleep_x nanosleep "libc.so.6"
//go:cgo_import_dynamic libc_sigfillset_x sigfillset "libc.so.6"
//go:cgo_import_dynamic libc_abort_x abort "libc.so.6"
//go:cgo_import_dynamic libc_dprintf_x dprintf "libc.so.6"
//go:cgo_import_dynamic libc_strerror_x strerror "libc.so.6"

// on amd64 we don't need the following lines - on 386 we do...
// anyway - with those lines the output is better (but doesn't matter) - without it on amd64 we get multiple DT_NEEDED with "libc.so.6" etc

//go:cgo_import_dynamic _ _ "libpthread.so.0"
//go:cgo_import_dynamic _ _ "libc.so.6"
