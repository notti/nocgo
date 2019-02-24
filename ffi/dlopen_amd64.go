package ffi

import _ "unsafe" // needed for go:linkname

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
