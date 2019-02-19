package testlib

import (
	"log"
	"os"
	"runtime"
	"testing"
	"unsafe"

	"github.com/notti/nocgo"
)

type emptySpec struct {

}

var emptyFunc nocgo.Spec

func TestEmpty(t *testing.T) {
	arg := &emptySpec{  }
	t.Log(emptyFunc)
	emptyFunc.Call(unsafe.Pointer(arg))
}

type int1Spec struct {
	ret int8 `nocgo:"ret"`
}

var int1Func nocgo.Spec

func TestInt1(t *testing.T) {
	arg := &int1Spec{  }
	t.Log(int1Func)
	int1Func.Call(unsafe.Pointer(arg))
	if arg.ret != 10 {
		t.Fatalf("Expected %v, but got %v\n", 10, arg.ret)
	}
}

type int2Spec struct {
	ret int8 `nocgo:"ret"`
}

var int2Func nocgo.Spec

func TestInt2(t *testing.T) {
	arg := &int2Spec{  }
	t.Log(int2Func)
	int2Func.Call(unsafe.Pointer(arg))
	if arg.ret != -10 {
		t.Fatalf("Expected %v, but got %v\n", -10, arg.ret)
	}
}

type int3Spec struct {
	ret uint8 `nocgo:"ret"`
}

var int3Func nocgo.Spec

func TestInt3(t *testing.T) {
	arg := &int3Spec{  }
	t.Log(int3Func)
	int3Func.Call(unsafe.Pointer(arg))
	if arg.ret != 10 {
		t.Fatalf("Expected %v, but got %v\n", 10, arg.ret)
	}
}

type int4Spec struct {
	ret uint8 `nocgo:"ret"`
}

var int4Func nocgo.Spec

func TestInt4(t *testing.T) {
	arg := &int4Spec{  }
	t.Log(int4Func)
	int4Func.Call(unsafe.Pointer(arg))
	if arg.ret != 246 {
		t.Fatalf("Expected %v, but got %v\n", 246, arg.ret)
	}
}

type float1Spec struct {
	ret float32 `nocgo:"ret"`
}

var float1Func nocgo.Spec

func TestFloat1(t *testing.T) {
	arg := &float1Spec{  }
	t.Log(float1Func)
	float1Func.Call(unsafe.Pointer(arg))
	if arg.ret != 10.5 {
		t.Fatalf("Expected %v, but got %v\n", 10.5, arg.ret)
	}
}

type float2Spec struct {
	ret float64 `nocgo:"ret"`
}

var float2Func nocgo.Spec

func TestFloat2(t *testing.T) {
	arg := &float2Spec{  }
	t.Log(float2Func)
	float2Func.Call(unsafe.Pointer(arg))
	if arg.ret != 10.5 {
		t.Fatalf("Expected %v, but got %v\n", 10.5, arg.ret)
	}
}



func TestMain(m *testing.M) {
	var lib string
	switch runtime.GOARCH {
	case "386":
		lib = "libcalltest32.so.1"
	case "amd64":
		lib = "libcalltest64.so.1"
	default:
		log.Fatalln("Unknown arch ", runtime.GOARCH)
	}

	l, err := nocgo.Open(lib)
	if err != nil {
		log.Fatal(err)
	}

	emptyFunc, err = l.Func("empty", emptySpec{})
	if err != nil {
		log.Fatal(err)
	}

	int1Func, err = l.Func("int1", int1Spec{})
	if err != nil {
		log.Fatal(err)
	}

	int2Func, err = l.Func("int2", int2Spec{})
	if err != nil {
		log.Fatal(err)
	}

	int3Func, err = l.Func("int3", int3Spec{})
	if err != nil {
		log.Fatal(err)
	}

	int4Func, err = l.Func("int4", int4Spec{})
	if err != nil {
		log.Fatal(err)
	}

	float1Func, err = l.Func("float1", float1Spec{})
	if err != nil {
		log.Fatal(err)
	}

	float2Func, err = l.Func("float2", float2Spec{})
	if err != nil {
		log.Fatal(err)
	}

		

	os.Exit(m.Run())
}
