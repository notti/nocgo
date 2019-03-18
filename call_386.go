package nocgo

import (
	"unsafe"
)

// 386 cdecl calling conventions: http://www.sco.com/developers/devspecs/abi386-4.pdf
// Pass everything on the stack in right to left order
// Return is in AX (and DX for 64 bit) or F0 for floats
// according to libffi clang might require the caller to properly (sign)extend stuff - so we do that
// structs are not supported for now (neither as argument nor as return value) - but this is not hard to do

type argtype uint16

const (
	type32     argtype = 0 // movl    64 bit
	type16     argtype = 1 // movw    16 bit
	type8      argtype = 2 // movb    8  bit
	typeDouble argtype = 3 // fld     64 bit (only return)
	typeFloat  argtype = 4 // movss   32 bit (only return)
	type64     argtype = 5 // 2x movl        (only return)
	typeUnused argtype = 0xFFFF
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
	ret     argument
}

func callWrapper()

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

	// on 386 we can't directly pass the arguments to the function :(
	// -> go aligns arguments like a struct and 386 aligns every argument to 4 byte boundaries
	for _, arg := range arguments {
		switch arg.size {
		case 8:
			spec.stack = append(spec.stack, argument{uint16(arg.offset), type32})
			spec.stack = append(spec.stack, argument{uint16(arg.offset + 4), type32})
		case 4:
			spec.stack = append(spec.stack, argument{uint16(arg.offset), type32})
		case 2:
			spec.stack = append(spec.stack, argument{uint16(arg.offset), type16})
		case 1:
			spec.stack = append(spec.stack, argument{uint16(arg.offset), type8})
		}
	}

	if ret.c != classVoid {
		var t argtype
		switch ret.c {
		case classInt, classUint:
			switch ret.size {
			case 8:
				t = type64
			case 4:
				t = type32
			case 2:
				t = type16
			case 1:
				t = type8
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
