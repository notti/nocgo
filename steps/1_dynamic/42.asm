; 42 in 64 bit (afzer tiny.asm)
BITS 64
GLOBAL __printf
TYPE __printf object
SECTION .bss
__printf:
    DD 0xDEADBEEF
    DD 0xDEADBEEF
.end:
size __printf __printf.end-__printf

SECTION .data
TYPE blah object
GLOBAL blah
blah db 'hello world', 0
.end:
size blah blah.end-blah

GLOBAL _start
SECTION .text
_start:
    mov     rdi, blah
    call    [__printf]

    mov     rax, 60
    mov     rdi, 42  
    syscall
TYPE _start FUNCTION