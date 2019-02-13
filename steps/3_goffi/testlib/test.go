package main

// int test_call(short a, int b, float c, double d);
import "C"
import "fmt"

func main() {
	fmt.Println(C.test_call(1, 2, 3, 4))
}
