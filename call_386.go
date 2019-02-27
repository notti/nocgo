package nocgo

import (
	"fmt"
	"reflect"
	"strings"
)

type argtype uint16

const (
	type32     argtype = 0 // movl              64 bit
	typeS16    argtype = 1 // movwlsx    signed 16 bit
	typeU16    argtype = 2 // movwlzx  unsigned 16 bit
	typeS8     argtype = 3 // movblsx    signed 8  bit
	typeU8     argtype = 4 // movblzx  unsigned 8  bit
	typeDouble argtype = 5 // fld             64 bit
	typeFloat  argtype = 6 // movss             32 bit
	type64     argtype = 7
	typeUnused argtype = 0xFFFF
)

type argument struct {
	offset uint16
	t      argtype
}

// Spec is the callspec needed to do the actuall call
type Spec struct {
	fn    uintptr
	base  uintptr
	stack []argument
	ret   argument
}

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
			switch f.Type.Kind() {
			case reflect.Slice:
				spec.ret = argument{uint16(f.Offset + sliceOffset), type32}
			case reflect.Int32, reflect.Uint32, reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
				spec.ret = argument{uint16(f.Offset), type32}
			case reflect.Int16:
				spec.ret = argument{uint16(f.Offset), typeS16}
			case reflect.Uint16:
				spec.ret = argument{uint16(f.Offset), typeU16}
			case reflect.Int8:
				spec.ret = argument{uint16(f.Offset), typeS8}
			case reflect.Uint8, reflect.Bool:
				spec.ret = argument{uint16(f.Offset), typeU8}
			case reflect.Float32:
				spec.ret = argument{uint16(f.Offset), typeFloat}
			case reflect.Float64:
				spec.ret = argument{uint16(f.Offset), typeDouble}
			case reflect.Uint64, reflect.Int64:
				spec.ret = argument{uint16(f.Offset), type64}
			default:
				panic("Unknown return Type")
			}
			continue
		}
		switch f.Type.Kind() {
		case reflect.Slice:
			spec.stack = append(spec.stack, argument{uint16(f.Offset + sliceOffset), type32})
		case reflect.Uint64, reflect.Int64, reflect.Float64:
			spec.stack = append(spec.stack, argument{uint16(f.Offset), type32})
			spec.stack = append(spec.stack, argument{uint16(f.Offset + 4), type32})
		case reflect.Int32, reflect.Uint32, reflect.Ptr, reflect.Uintptr, reflect.Float32:
			spec.stack = append(spec.stack, argument{uint16(f.Offset), type32})
		case reflect.Int16:
			spec.stack = append(spec.stack, argument{uint16(f.Offset), typeS16})
		case reflect.Uint16:
			spec.stack = append(spec.stack, argument{uint16(f.Offset), typeU16})
		case reflect.Int8:
			spec.stack = append(spec.stack, argument{uint16(f.Offset), typeS8})
		case reflect.Uint8, reflect.Bool:
			spec.stack = append(spec.stack, argument{uint16(f.Offset), typeU8})
		default:
			fmt.Println(f.Type.Kind())
			panic("Unknown type")
		}
	}
	return spec, nil
}
