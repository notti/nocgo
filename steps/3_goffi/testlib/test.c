#include <stdio.h>

int test_call(short a, int b, float c, double d) {
	printf("In C: %d %d %f %f\n", a, b, c, d);
	return a+b+c+d;
}
