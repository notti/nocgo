package main

// int test_call(unsigned short i1, int i2, float f1, double f2, int i3, int i4, int i5, int i6, int i7, int i8, char i9);
import "C"

import "fmt"

func main() {
	fmt.Println(C.test_call(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, -11))
}
