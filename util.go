package nocgo

import "unsafe"

// MakeCString converts the given string to a null terminated C byte-slice/char array
func MakeCString(s string) []byte {
	return append([]byte(s), 0)
}

// MakeGoStringFromPointer converts the given pointer to a null terminated C string to a go string
func MakeGoStringFromPointer(s uintptr) string {
	if s == 0 {
		return ""
	}
	bval := (*[1 << 30]byte)(unsafe.Pointer(s))
	return MakeGoStringFromSlice(bval[:])
}

// MakeGoStringFromSlice converts the given byte slice containing a null terminated C string to a go string
func MakeGoStringFromSlice(s []byte) string {
	if len(s) == 0 {
		return ""
	}
	for i := range s {
		if s[i] == 0 {
			return string(s[:i])
		}
	}
	return string(s[:])
}
