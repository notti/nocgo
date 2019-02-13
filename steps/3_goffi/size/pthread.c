#include <pthread.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	printf("%d %d\n", sizeof(pthread_attr_t), sizeof(size_t));
	return 0;
}
