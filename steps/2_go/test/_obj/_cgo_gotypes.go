// Code generated by cmd/cgo; DO NOT EDIT.

package main

import "unsafe"

import _ "runtime/cgo"

import "syscall"

var _ syscall.Errno
func _Cgo_ptr(ptr unsafe.Pointer) unsafe.Pointer { return ptr }

//go:linkname _Cgo_always_false runtime.cgoAlwaysFalse
var _Cgo_always_false bool
//go:linkname _Cgo_use runtime.cgoUse
func _Cgo_use(interface{})
type _Ctype_FILE = _Ctype_struct__IO_FILE

type _Ctype__IO_lock_t = _Ctype_void

type _Ctype___off64_t = _Ctype_long

type _Ctype___off_t = _Ctype_long

type _Ctype_char int8

type _Ctype_int int32

type _Ctype_long int64

type _Ctype_schar int8

type _Ctype_size_t = _Ctype_ulong

type _Ctype_struct__IO_FILE struct {
	_flags		_Ctype_int
	_IO_read_ptr	*_Ctype_char
	_IO_read_end	*_Ctype_char
	_IO_read_base	*_Ctype_char
	_IO_write_base	*_Ctype_char
	_IO_write_ptr	*_Ctype_char
	_IO_write_end	*_Ctype_char
	_IO_buf_base	*_Ctype_char
	_IO_buf_end	*_Ctype_char
	_IO_save_base	*_Ctype_char
	_IO_backup_base	*_Ctype_char
	_IO_save_end	*_Ctype_char
	_markers	*_Ctype_struct__IO_marker
	_chain		*_Ctype_struct__IO_FILE
	_fileno		_Ctype_int
	_flags2		_Ctype_int
	_old_offset	_Ctype___off_t
	_cur_column	_Ctype_ushort
	_vtable_offset	_Ctype_schar
	_shortbuf	[1]_Ctype_char
	_lock		unsafe.Pointer
	_offset		_Ctype___off64_t
	_codecvt	*_Ctype_struct__IO_codecvt
	_wide_data	*_Ctype_struct__IO_wide_data
	_freeres_list	*_Ctype_struct__IO_FILE
	_freeres_buf	unsafe.Pointer
	__pad5		_Ctype_size_t
	_mode		_Ctype_int
	_unused2	[20]_Ctype_char
}

type _Ctype_struct__IO_codecvt struct{}

type _Ctype_struct__IO_marker struct{}

type _Ctype_struct__IO_wide_data struct{}

type _Ctype_ulong uint64

type _Ctype_ushort uint16

type _Ctype_void [0]byte

//go:linkname _cgo_runtime_cgocall runtime.cgocall
func _cgo_runtime_cgocall(unsafe.Pointer, uintptr) int32

//go:linkname _cgo_runtime_cgocallback runtime.cgocallback
func _cgo_runtime_cgocallback(unsafe.Pointer, unsafe.Pointer, uintptr, uintptr)

//go:linkname _cgoCheckPointer runtime.cgoCheckPointer
func _cgoCheckPointer(interface{}, ...interface{})

//go:linkname _cgoCheckResult runtime.cgoCheckResult
func _cgoCheckResult(interface{})

//go:cgo_import_static _cgo_4b63218df55f_Cfunc_fputs
//go:linkname __cgofn__cgo_4b63218df55f_Cfunc_fputs _cgo_4b63218df55f_Cfunc_fputs
var __cgofn__cgo_4b63218df55f_Cfunc_fputs byte
var _cgo_4b63218df55f_Cfunc_fputs = unsafe.Pointer(&__cgofn__cgo_4b63218df55f_Cfunc_fputs)

//go:cgo_unsafe_args
func _Cfunc_fputs(p0 *_Ctype_char, p1 *_Ctype_struct__IO_FILE) (r1 _Ctype_int) {
	_cgo_runtime_cgocall(_cgo_4b63218df55f_Cfunc_fputs, uintptr(unsafe.Pointer(&p0)))
	if _Cgo_always_false {
		_Cgo_use(p0)
		_Cgo_use(p1)
	}
	return
}
//go:cgo_import_static _cgo_4b63218df55f_Cfunc_putc
//go:linkname __cgofn__cgo_4b63218df55f_Cfunc_putc _cgo_4b63218df55f_Cfunc_putc
var __cgofn__cgo_4b63218df55f_Cfunc_putc byte
var _cgo_4b63218df55f_Cfunc_putc = unsafe.Pointer(&__cgofn__cgo_4b63218df55f_Cfunc_putc)

//go:cgo_unsafe_args
func _Cfunc_putc(p0 _Ctype_int, p1 *_Ctype_struct__IO_FILE) (r1 _Ctype_int) {
	_cgo_runtime_cgocall(_cgo_4b63218df55f_Cfunc_putc, uintptr(unsafe.Pointer(&p0)))
	if _Cgo_always_false {
		_Cgo_use(p0)
		_Cgo_use(p1)
	}
	return
}
//go:cgo_import_static _cgo_4b63218df55f_Cfunc_puts
//go:linkname __cgofn__cgo_4b63218df55f_Cfunc_puts _cgo_4b63218df55f_Cfunc_puts
var __cgofn__cgo_4b63218df55f_Cfunc_puts byte
var _cgo_4b63218df55f_Cfunc_puts = unsafe.Pointer(&__cgofn__cgo_4b63218df55f_Cfunc_puts)

//go:cgo_unsafe_args
func _Cfunc_puts(p0 *_Ctype_char) (r1 _Ctype_int) {
	_cgo_runtime_cgocall(_cgo_4b63218df55f_Cfunc_puts, uintptr(unsafe.Pointer(&p0)))
	if _Cgo_always_false {
		_Cgo_use(p0)
	}
	return
}
//go:cgo_import_static _cgo_4b63218df55f_Cmacro_stdout
//go:linkname __cgofn__cgo_4b63218df55f_Cmacro_stdout _cgo_4b63218df55f_Cmacro_stdout
var __cgofn__cgo_4b63218df55f_Cmacro_stdout byte
var _cgo_4b63218df55f_Cmacro_stdout = unsafe.Pointer(&__cgofn__cgo_4b63218df55f_Cmacro_stdout)

//go:cgo_unsafe_args
func _Cmacro_stdout() (r1 *_Ctype_struct__IO_FILE) {
	_cgo_runtime_cgocall(_cgo_4b63218df55f_Cmacro_stdout, uintptr(unsafe.Pointer(&r1)))
	if _Cgo_always_false {
	}
	return
}
//go:cgo_import_static _cgo_4b63218df55f_Cfunc_strcat
//go:linkname __cgofn__cgo_4b63218df55f_Cfunc_strcat _cgo_4b63218df55f_Cfunc_strcat
var __cgofn__cgo_4b63218df55f_Cfunc_strcat byte
var _cgo_4b63218df55f_Cfunc_strcat = unsafe.Pointer(&__cgofn__cgo_4b63218df55f_Cfunc_strcat)

//go:cgo_unsafe_args
func _Cfunc_strcat(p0 *_Ctype_char, p1 *_Ctype_char) (r1 *_Ctype_char) {
	_cgo_runtime_cgocall(_cgo_4b63218df55f_Cfunc_strcat, uintptr(unsafe.Pointer(&p0)))
	if _Cgo_always_false {
		_Cgo_use(p0)
		_Cgo_use(p1)
	}
	return
}
//go:cgo_import_static _cgo_4b63218df55f_Cfunc_test_cb
//go:linkname __cgofn__cgo_4b63218df55f_Cfunc_test_cb _cgo_4b63218df55f_Cfunc_test_cb
var __cgofn__cgo_4b63218df55f_Cfunc_test_cb byte
var _cgo_4b63218df55f_Cfunc_test_cb = unsafe.Pointer(&__cgofn__cgo_4b63218df55f_Cfunc_test_cb)

//go:cgo_unsafe_args
func _Cfunc_test_cb() (r1 _Ctype_int) {
	_cgo_runtime_cgocall(_cgo_4b63218df55f_Cfunc_test_cb, uintptr(unsafe.Pointer(&r1)))
	if _Cgo_always_false {
	}
	return
}
//go:cgo_export_dynamic cb
//go:linkname _cgoexp_4b63218df55f_cb _cgoexp_4b63218df55f_cb
//go:cgo_export_static _cgoexp_4b63218df55f_cb
//go:nosplit
//go:norace
func _cgoexp_4b63218df55f_cb(a unsafe.Pointer, n int32, ctxt uintptr) {
	fn := _cgoexpwrap_4b63218df55f_cb
	_cgo_runtime_cgocallback(**(**unsafe.Pointer)(unsafe.Pointer(&fn)), a, uintptr(n), ctxt);
}

func _cgoexpwrap_4b63218df55f_cb(p0 int) (r0 int) {
	return cb(p0)
}
