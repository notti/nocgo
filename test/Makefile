all:
	@go run build.go
	@make -C testlib
	@make -C ffi

verbose:
	@make all TEST="-test.v"

bench:
	@make all TEST="-test.bench . -test.benchmem"


.PHONY: all
