  ; 42 in 64 bit (afzer tiny.asm)
  BITS 64
  GLOBAL _start
  SECTION .text
  _start:
                mov     rax, 60
                mov     rdi, 42  
                syscall
