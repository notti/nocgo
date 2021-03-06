// +build !cgo

package dlopen

// we have to use the 3 argument format here :( - 2 argument format is only allowed from inside cgo

//go:cgo_import_dynamic libc_dlopen_x dlopen "libc.so.7"
//go:cgo_import_dynamic libc_dlclose_x dlclose "libc.so.7"
//go:cgo_import_dynamic libc_dlsym_x dlsym "libc.so.7"
//go:cgo_import_dynamic libc_dlerror_x dlerror "libc.so.7"

// on amd64 we don't need the following line - on 386 we do...
// anyway - with those lines the output is better (but doesn't matter) - without it on amd64 we get multiple DT_NEEDED with "libc.so.6" etc

//go:cgo_import_dynamic _ _ "libc.so.7"
