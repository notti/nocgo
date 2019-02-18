all: main32 main64

main64: main.go
	go build -o ../main64

main32: main.go
	GOARCH=386 go build -o ../main32