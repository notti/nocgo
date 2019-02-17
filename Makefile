main:
	rm -f main
	$(MAKE) -C test 
	$(MAKE) -C testlib
	go run relink/relink.go main
	LD_LIBRARY_PATH=testlib ./main a b c

.PHONY: main
