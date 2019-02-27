platform_id = $$( uname -s )
GO_GCFLAGS = $$( \
    case $(platform_id) in \
        (FreeBSD) echo -gcflags=github.com/notti/nocgo/fakecgo=-std ;; \
    esac )

all:
	CGO_ENABLED=0 go build $(GO_GCFLAGS)
