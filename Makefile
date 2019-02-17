main: relink.go
	rm -f main
	$(MAKE) -C test 
	$(MAKE) -C testlib
	go run relink.go main
	LD_LIBRARY_PATH=testlib ./main a b c

.PHONY: main
