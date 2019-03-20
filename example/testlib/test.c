#include <stdio.h>

int value;

int test_call(unsigned short i1, int i2, float f1, double f2, int i3, int i4, int i5, int i6, int i7, int i8, char i9) {
	printf("In C: %d %d %f %f\n", i1, i2, f1, f2);
	return i1+i2+f1+f2+i3+i4+i5+i6+i7+i8+i9;
}

void print_value() {
	printf("value: %d\n", value);
}

int test_cb(int (callback(int))) {
	int res;
	printf("in C!\n");
	res = callback(10);
	printf("in C again: %d\n", res);
	return res*2;
}
