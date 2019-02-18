package nocgo

import (
	"reflect"
	"strings"
)

type argtype uint16

const (
	type64     argtype = 0 // movq              64 bit
	typeS32    argtype = 1 // movlqsx    signed 32 bit
	typeU32    argtype = 2 // movlqzx  unsigned 32 bit
	typeS16    argtype = 3 // movwqsx    signed 16 bit
	typeU16    argtype = 4 // movwqzx  unsigned 16 bit
	typeS8     argtype = 5 // movbqsx    signed 8  bit
	typeU8     argtype = 6 // movbqzx  unsigned 8  bit
	typeDouble argtype = 7 // movsd             64 bit
	typeFloat  argtype = 8 // movss             32 bit
	typeUnused argtype = 0xFFFF
)

type argument struct {
	offset uint16
	t      argtype
}

// Spec is the callspec needed to do the actuall call
type Spec struct {
	fn      uintptr
	base    uintptr
	stack   []argument
	intargs [6]argument
	xmmargs [8]argument
	ret     argument
	rax     uint8
}

func fieldToOffset(k reflect.StructField) (argument, bool) {
	switch k.Type.Kind() {
	case reflect.Slice:
		return argument{uint16(k.Offset + sliceOffset), type64}, false
	case reflect.Int, reflect.Uint, reflect.Uint64, reflect.Int64, reflect.Ptr, reflect.Uintptr:
		return argument{uint16(k.Offset), type64}, false
	case reflect.Int32:
		return argument{uint16(k.Offset), typeS32}, false
	case reflect.Uint32:
		return argument{uint16(k.Offset), typeU32}, false
	case reflect.Int16:
		return argument{uint16(k.Offset), typeS16}, false
	case reflect.Uint16:
		return argument{uint16(k.Offset), typeU16}, false
	case reflect.Int8:
		return argument{uint16(k.Offset), typeS8}, false
	case reflect.Uint8, reflect.Bool:
		return argument{uint16(k.Offset), typeU8}, false
	case reflect.Float32:
		return argument{uint16(k.Offset), typeFloat}, true
	case reflect.Float64:
		return argument{uint16(k.Offset), typeDouble}, true
	}
	panic("Unknown Type")
}

// FIXME: we don't support stuff > 64 bit

// makeSpec builds a call specification for the given arguments
func makeSpec(fn uintptr, args interface{}) (Spec, error) {
	v := reflect.ValueOf(args)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	var spec Spec

	spec.fn = fn

	spec.ret.t = typeUnused

	haveRet := false

	intreg := 0
	xmmreg := 0

ARGS:
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tags := strings.Split(f.Tag.Get("nocgo"), ",")
		ret := false
		for _, tag := range tags {
			if tag == "ignore" {
				continue ARGS
			}
			if tag == "ret" {
				if haveRet == true {
					panic("Only one return argument allowed")
				}
				ret = true
				haveRet = true
				continue
			}
		}
		if ret {
			off, _ := fieldToOffset(f)
			spec.ret = off
			// FIXME ret1/xmmret1! - only needed for types > 64 bit
			continue
		}
		off, xmm := fieldToOffset(f)
		if xmm {
			if xmmreg < 8 {
				spec.xmmargs[xmmreg] = off
				xmmreg++
			} else {
				spec.stack = append(spec.stack, off)
			}
		} else {
			if intreg < 6 {
				spec.intargs[intreg] = off
				intreg++
			} else {
				spec.stack = append(spec.stack, off)
			}
		}
	}
	for i := intreg; i < 6; i++ {
		spec.intargs[i].t = typeUnused
	}
	for i := xmmreg; i < 8; i++ {
		spec.xmmargs[i].t = typeUnused
	}
	spec.rax = uint8(xmmreg)
	return spec, nil
}
