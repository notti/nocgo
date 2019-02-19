package testlib

import (
	"testing"
)

func TestEmpty(t *testing.T) {
	empty()
}

func TestInt1(t *testing.T) {
	ret := int1()
	if ret != 10 {
		t.Fatalf("Expected %v, but got %v\n", 10, ret)
	}
}

func TestInt2(t *testing.T) {
	ret := int2()
	if ret != -10 {
		t.Fatalf("Expected %v, but got %v\n", -10, ret)
	}
}

func TestInt3(t *testing.T) {
	ret := int3()
	if ret != 10 {
		t.Fatalf("Expected %v, but got %v\n", 10, ret)
	}
}

func TestInt4(t *testing.T) {
	ret := int4()
	if ret != 246 {
		t.Fatalf("Expected %v, but got %v\n", 246, ret)
	}
}

func TestFloat1(t *testing.T) {
	ret := float1()
	if ret != 10.5 {
		t.Fatalf("Expected %v, but got %v\n", 10.5, ret)
	}
}

func TestFloat2(t *testing.T) {
	ret := float2()
	if ret != 10.5 {
		t.Fatalf("Expected %v, but got %v\n", 10.5, ret)
	}
}


