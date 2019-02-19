package testlib

// void empty();
// char int1();
// char int2();
// unsigned char int3();
// unsigned char int4();
// float float1();
// double float2();
import "C"

func empty()  {
	C.empty()
}

func int1() int8 {
	return int8(C.int1())
}

func int2() int8 {
	return int8(C.int2())
}

func int3() uint8 {
	return uint8(C.int3())
}

func int4() uint8 {
	return uint8(C.int4())
}

func float1() float32 {
	return float32(C.float1())
}

func float2() float64 {
	return float64(C.float2())
}



