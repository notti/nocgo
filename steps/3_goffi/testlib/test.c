#include <stdio.h>

int test_call(short i1, int i2, float f1, double f2, int i3, int i4, int i5, int i6, int i7, int i8, short i9) {
	printf("In C: %d %d %f %f\n", i1, i1, f1, f2);
	return i1+i2+f1+f2+i3+i4+i5+i6+i7+i8+i9;
}
