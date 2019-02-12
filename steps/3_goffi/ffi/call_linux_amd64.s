#include "textflag.h"

/*
    Frame layout:
/*
	0:   fn
	8:   base
	16:  stack
	24:  slicelen
	32:  slicecap
	40:  regarg0 rdi
	44:  regarg1 rsi
	48:  regarg2 rdx
	52:  regarg3 rcx
	56:  regarg4 r8
	60:  regarg5 r9
	64:  xmmarg0 xmm2
	68:  xmmarg1 xmm3
	72:  xmmarg2 xmm4
	76:  xmmarg3 xmm5
	80:  xmmarg4 xmm6
	84:  xmmarg5 xmm7
	88:  ret0    rax
	92:  ret1    rdx
	96:  xmmret0 xmm0
	100: xmmret1 xmm1
	104: rax     rax

    int:
    0: movq     64bit
    1: movlqsx    signed 32 bit
    2: movlqzx  unsigned 32 bit
    3: movwqsx    signed 16 bit
    4: movwqzx  unsigned 16 bit
    5: movbqsx    signed 8  bit
    6: movbqzx  unsigned 8  bit

    float:
    0: movsd             64 bit
    1: movss             32 bit
*/

#define LOADREG(off, target) \
    MOVLQSX off(R12), AX \
    TESTQ AX, AX \
    JS prepared \
    MOVWQZX AX, R11 \
    SHRQ $16, AX \
    ADDQ R13, R11 \
    CMPB AX, $0 \
    JNE 3(PC) \
    MOVQ 0(R11), target \
    JMP 20(PC) \
    CMPB AX, $1 \
    JNE 3(PC) \
    MOVLQSX 0(R11), target \
    JMP 18(PC) \
    CMPB AX, $2 \
    JNE 3(PC) \
    MOVLQZX 0(R11), target \
    JMP 14(PC) \
    CMPB AX, $3 \
    JNE 3(PC) \
    MOVWQSX 0(R11), target \
    JMP 10(PC) \
    CMPB AX, $4 \
    JNE 3(PC) \
    MOVWQZX 0(R11), target \
    JMP 6(PC) \
    CMPB AX, $5 \
    JNE 3(PC) \
    MOVBQSX 0(R11), target \
    JMP 2(PC) \
    MOVBQZX 0(R11), target

// func asmcall3()
TEXT ·asmcall3(SB),NOSPLIT,$0
    MOVQ DI, R12      // FRAME (preserved)
    MOVQ 8(R12), R13  // base

    // load register arguments
    LOADREG(40, DI)
    LOADREG(44, SI)
    LOADREG(48, DX)
    LOADREG(52, CX)
    LOADREG(56, R8)
    LOADREG(60, R9)

prepared:

/*
    float: movss
    double: movsd
*/

    // load number of vector registers
    MOVQ 64(R12), AX

    // do the actuall call
    CALL (R12)

    // FIXME: correct return type
    // store ret0
    MOVLQSX 88(R12), BX
    TESTQ BX, BX
    JS DONE
    MOVWQZX BX, R11
    ADDQ R13, R11
    MOVQ AX, (R11)

DONE:
    RET
