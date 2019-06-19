platform_id = $$( uname -s )
GO_GCFLAGS = $$( \
    case $(platform_id) in \
        (FreeBSD) echo -gcflags=github.com/notti/nocgo/fakecgo=-std ;; \
    esac )

all: main32 main32C main64 main64C testlib

testlib: 
	$(MAKE) -C testlib

main32: main.go testlib
	rm -f main32
	-CGO_ENABLED=0 GOARCH=386 go build $(GO_GCFLAGS) -o main32 main.go && LD_LIBRARY_PATH=testlib ./main32 a b c
	
main32C: main.go testlib
	rm -f main32C
	-CGO_ENABLED=1 GOARCH=386 go build $(GO_GCFLAGS) -o main32C main.go && LD_LIBRARY_PATH=testlib ./main32C a b c

main64: main.go testlib
	rm -f main64
	-CGO_ENABLED=0 go build $(GO_GCFLAGS) -o main64 main.go && LD_LIBRARY_PATH=testlib ./main64 a b c

main64C: main.go testlib
	rm -f main64C
	-CGO_ENABLED=0 go build $(GO_GCFLAGS) -o main64C main.go && LD_LIBRARY_PATH=testlib ./main64C a b c

clean:
	rm -rf main32 main64 main32C main64C
	$(MAKE) -C testlib clean

.PHONY: main clean main32 main64 testlib
