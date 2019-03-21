package nocgo

import (
	"fmt"
	"unsafe"
)

// amd64 cdecl calling conventions: https://www.uclibc.org/docs/psABI-x86_64.pdf
//   - Align the stack (should be 32byte aligned before the function is called - 16 byte is enough if we don't pass 256bit integers)
//   - Pass first integer arguments in DI, SI, DX, CX, R8, R9
//   - Pass first float arguments in X0-X7
//   - Pass rest on the stack
//   - Pass number of used float registers in AX
// Return is in AX or X0 for floats
// according to libffi clang might require the caller to properly (sign)extend stuff in registers - so we do that
// structs are not supported for now (neither as argument nor as return value) - but this is not hard to do

type argtype uint16

const (
	type64       argtype = 0 // movq              64 bit
	typeS32      argtype = 1 // movlqsx    signed 32 bit
	typeU32      argtype = 2 // movlqzx  unsigned 32 bit
	typeS16      argtype = 3 // movwqsx    signed 16 bit
	typeU16      argtype = 4 // movwqzx  unsigned 16 bit
	typeS8       argtype = 5 // movbqsx    signed 8  bit
	typeU8       argtype = 6 // movbqzx  unsigned 8  bit
	typeDouble   argtype = 7 // movsd             64 bit
	typeFloat    argtype = 8 // movss             32 bit
	typeCallback argtype = 9
	typeUnused   argtype = 0xFFFF
)

type argument struct {
	offset uint16
	t      argtype
}

// spec a wrapper specifcation with instructions on how to place arguments into registers/stack
type spec struct {
	wrapper uintptr // pointer to callWrapper()
	fn      uintptr // pointer to the C-function
	stack   []argument
	intargs [6]argument
	xmmargs [8]argument
	ret     argument
	rax     uint8
}

// FIXME: we don't support stuff > 64 bit

func callWrapper()

type callbackArgs struct {
	bp      uintptr
	intargs [6]uintptr
	xmmargs [8]uintptr
	ax      uintptr
	which   uintptr
	spec    *spec
}

func testCallback(args *callbackArgs) {
	/*
		TODO:
		-build frame
		-call function
		-set return value
	*/
	fmt.Printf("got: %#v\n%v\n", args, args.spec)
	args.ax = args.intargs[0] * 2
}

// makeSpec builds a call specification for the given arguments
func makeSpec(fn uintptr, fun interface{}) error {
	fptr, arguments, ret, err := stackFields(fun)
	if err != nil {
		return err
	}

	spec := new(spec)
	spec.wrapper = funcPC(callWrapper)
	spec.fn = fn
	spec.ret.t = typeUnused

	intreg := 0
	xmmreg := 0
	cbnum := 0

	for _, arg := range arguments {
		var t argtype
		switch arg.c {
		case classInt, classUint, classCallback:
			switch {
			case arg.size == 8:
				t = type64
			case arg.size == 4:
				if arg.c == classInt {
					t = typeS32
				} else {
					t = typeU32
				}
			case arg.size == 2:
				if arg.c == classInt {
					t = typeS16
				} else {
					t = typeU16
				}
			case arg.size == 1:
				if arg.c == classInt {
					t = typeS8
				} else {
					t = typeU8
				}
			}
			if arg.c == classCallback {
				t = typeCallback
				arg.offset = cbnum
				cbnum++
			}
			if intreg < 6 {
				spec.intargs[intreg] = argument{uint16(arg.offset), t}
				intreg++
			} else {
				switch t {
				case typeS32:
					t = typeU32
				case typeS16:
					t = typeU16
				case typeS8:
					t = typeU8
				}
				spec.stack = append(spec.stack, argument{uint16(arg.offset), t})
			}
		case classFloat:
			switch {
			case arg.size == 8:
				t = typeDouble
			case arg.size == 4:
				t = typeFloat
			}
			if xmmreg < 8 {
				spec.xmmargs[xmmreg] = argument{uint16(arg.offset), t}
				xmmreg++
			} else {
				switch t {
				case typeDouble:
					t = type64
				case typeFloat:
					t = typeU32
				}
				spec.stack = append(spec.stack, argument{uint16(arg.offset), t})
			}
		}
	}

	// check cbnum!

	spec.rax = uint8(xmmreg)
	for i := intreg; i < 6; i++ {
		spec.intargs[i].t = typeUnused
	}
	for i := xmmreg; i < 8; i++ {
		spec.xmmargs[i].t = typeUnused
	}

	if ret.c != classVoid {
		var t argtype
		switch ret.c {
		case classInt, classUint:
			switch ret.size {
			case 8:
				t = type64
			case 4:
				t = typeU32
			case 2:
				t = typeU16
			case 1:
				t = typeU8
			}
		case classFloat:
			switch ret.size {
			case 8:
				t = typeDouble
			case 4:
				t = typeFloat
			}
		}
		spec.ret.t = t
		spec.ret.offset = uint16(ret.offset)
	}

	*(*unsafe.Pointer)(fptr) = unsafe.Pointer(spec)

	return nil
}
