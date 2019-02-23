#include <pthread.h>
#include <stdio.h>
#include <time.h>
#include <signal.h>

int main(int argc, char *argv[]) {
	printf("%d %d %d %d %d\n", sizeof(struct timespec), sizeof(sigset_t), sizeof(pthread_t), sizeof(pthread_attr_t), sizeof(size_t));
	return 0;
}
