ffi.test:
	@rm -f ffi.test32 ffi.test64
	@go test -c -o ffi.test64
	@GOARCH=386 go test -c -o ffi.test32
	@go run ../../relink/relink.go ffi.test64
	@go run ../../relink/relink.go ffi.test32
	LD_LIBRARY_PATH="../testlib" ./ffi.test32 $(TEST)
	LD_LIBRARY_PATH="../testlib" ./ffi.test64 $(TEST)

.PHONY: ffi.test
