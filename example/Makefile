platform_id = $$( uname -s )
GO_GCFLAGS = $$( \
    case $(platform_id) in \
        (FreeBSD) echo -gcflags=github.com/notti/nocgo/fakecgo=-std ;; \
    esac )

all: main32 main64

testlib:
	$(MAKE) -C testlib

main32: main.go testlib
	rm -f main32
	-GOARCH=386 go build $(GO_GCFLAGS) -o main32 main.go && LD_LIBRARY_PATH=testlib ./main32 a b c

main64: main.go testlib
	rm -f main64
	-go build $(GO_GCFLAGS) -o main64 main.go && LD_LIBRARY_PATH=testlib ./main64 a b c

clean:
	rm -rf main32 main64
	$(MAKE) -C testlib clean

.PHONY: main clean main32 main64
