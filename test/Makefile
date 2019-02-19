all:
	@go run build.go
	@make -C testlib
	@make -C nocgo 

verbose:
	@make all TEST="-test.v"

.PHONY: all
