nocgo.test:
	@rm -f nocgo.test32 nocgo.test64
	@go test -c -o nocgo.test64
	@GOARCH=386 go test -c -o nocgo.test32
	@go run ../../relink/relink.go nocgo.test64
	@go run ../../relink/relink.go nocgo.test32
	LD_LIBRARY_PATH="../testlib" ./nocgo.test32 $(TEST)
	LD_LIBRARY_PATH="../testlib" ./nocgo.test64 $(TEST)

.PHONY: nocgo.test
