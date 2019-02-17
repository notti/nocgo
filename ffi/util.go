package ffi

import "unsafe"

// MakeCString converts the given string to a null terminated C string
func MakeCString(s string) []byte {
	return append([]byte(s), 0)
}

// MakeGoString converts the given pointer to a null terminated C string to a go string
func MakeGoString(s uintptr) string {
	if s == 0 {
		return ""
	}
	bval := (*[1 << 30]byte)(unsafe.Pointer(s))
	for i := range bval {
		if bval[i] == 0 {
			return string(bval[:i])
		}
	}
	return string(bval[:])
}
