a.out: 42.asm
	yasm -f elf -m amd64 42.asm
	ld -s 42.o
simple_debug: simple.c
	gcc -Wall -nostdlib -g -o simple_debug simple.c
simple: simple.c
	gcc -Wall -nostdlib -g -o simple simple.c
