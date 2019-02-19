nocgo.test:
	@rm -f nocgo.test
	@go test -c
	@go run ../../relink/relink.go nocgo.test
	@LD_LIBRARY_PATH="../testlib" ./nocgo.test $(TEST)

.PHONY: nocgo.test
