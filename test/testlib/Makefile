all: libcalltest32.so.1 libcalltest64.so.1 test

libcalltest64.so.1: test.c
	@$(CC) -g -shared -fPIC -Wl,-soname,libcalltest64.so.1 -o libcalltest64.so.1 test.c -lc

libcalltest32.so.1: test.c
	@$(CC) -g -shared -m32 -fPIC -Wl,-soname,libcalltest32.so.1 -o libcalltest32.so.1 test.c -lc

test:
	go test $(TEST)

.PHONY: test
