package main

/*
#include <stdio.h>
extern int cb(int);
int test_cb() {
	int res;
	printf("in C!\n");
	res = cb(10);
	printf("in C again: %d\n", res);
	return res*2;
}
*/
import "C"
