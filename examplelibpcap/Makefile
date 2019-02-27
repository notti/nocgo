platform_id = $$( uname -s )
GO_GCFLAGS = $$( \
	case $(platform_id) in \
		(FreeBSD) echo -gcflags=github.com/notti/nocgo/fakecgo=-std ;; \
	esac )

all:
	@rm -rf pcapnocgo pcapcgo
	CGO_ENABLED=0 go build $(GO_GCFLAGS) -o pcapnocgo
	CGO_ENABLED=1 go build -o pcapcgo

clean:
	rm -rf pcapnocgo pcapcgo

.PHONY: clean
