dyn: prog rewrite.go
	$(MAKE) -C prog
	$(MAKE) -C testlib
	go run rewrite.go
	LD_LIBRARY_PATH=testlib ./dyn a b c

libcalltest.so.1: test.c
	$(CC) -shared -fPIC -Wl,-soname,libcalltest.so.1 -o libcalltest.so.1 test.c -lc

.PHONY: dyn
