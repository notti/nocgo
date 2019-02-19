ffi.test:
	@rm -f ffi.test
	@go test -c
	@go run ../../relink/relink.go ffi.test
	@LD_LIBRARY_PATH="../testlib" ./ffi.test $(TEST)

.PHONY: ffi.test
