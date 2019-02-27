platform_id = $$( uname -s )
GO_GCFLAGS = $$( \
    case $(platform_id) in \
        (FreeBSD) echo -gcflags=github.com/notti/nocgo/fakecgo=-std ;; \
    esac )

all: nocgo.test32 nocgo.test64

nocgo.test32:
	@rm -f nocgo.test32
	-@CGO_ENABLED=0 GOARCH=386 go test -c -o nocgo.test32 $(GO_GCFLAGS) && LD_LIBRARY_PATH="../testlib" ./nocgo.test32 $(TEST)

nocgo.test64:
	@rm -f nocgo.test64
	-@CGO_ENABLED=0 go test -c -o nocgo.test64 $(GO_GCFLAGS) && LD_LIBRARY_PATH="../testlib" ./nocgo.test64 $(TEST)

clean:
	rm -f nocgo.test32 nocgo.test64

.PHONY: nocgo.test32 nocgo.test64 clean
