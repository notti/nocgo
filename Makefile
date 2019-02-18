main: relink.go
	rm -f main64 main32
	$(MAKE) -C test
	$(MAKE) -C testlib
	go run relink.go main32
	go run relink.go main64
	LD_LIBRARY_PATH=testlib ./main32 a b c
	LD_LIBRARY_PATH=testlib ./main64 a b c

.PHONY: main
