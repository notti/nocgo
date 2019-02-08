a.out: 42.asm
	yasm -f elf -m amd64 -g dwarf2 42.asm
	ld 42.o
dyn: a.out rewrite.go
	go run rewrite.go
