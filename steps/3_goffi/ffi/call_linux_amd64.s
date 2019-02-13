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
	64:  xmmarg0 xmm0
	68:  xmmarg1 xmm1
	72:  xmmarg2 xmm2
	76:  xmmarg3 xmm3
	80:  xmmarg4 xmm4
	84:  xmmarg5 xmm5
	88:  xmmarg6 xmm6
	92:  xmmarg7 xmm7
	96:  ret0    rax
	100: ret1    rdx
	104: xmmret0 xmm0
	108: xmmret1 xmm1
	112: rax     rax

    int:
    0: movq     64bit
    1: movlqsx    signed 32 bit
    2: movlqzx  unsigned 32 bit
    3: movwqsx    signed 16 bit
    4: movwqzx  unsigned 16 bit
    5: movbqsx    signed 8  bit
    6: movbqzx  unsigned 8  bit

    float:
    7: movsd             64 bit
    8: movss             32 bit
*/

#define LOADREG(off, target) \
    MOVLQSX off(R12), AX \
    TESTQ AX, AX \
    JS xmm \
    MOVWQZX AX, R11 \
    SHRL $16, AX \
    ADDQ R13, R11 \
    CMPB AX, $0 \
    JNE 3(PC) \
    MOVQ 0(R11), target \ // 64bit
    JMP 20(PC) \
    CMPB AX, $1 \
    JNE 3(PC) \
    MOVLQSX 0(R11), target \ // signed 32 bit
    JMP 18(PC) \
    CMPB AX, $2 \
    JNE 3(PC) \
    MOVLQZX 0(R11), target \ // unsigned 32 bit
    JMP 14(PC) \
    CMPB AX, $3 \
    JNE 3(PC) \
    MOVWQSX 0(R11), target \ // signed 16 bit
    JMP 10(PC) \
    CMPB AX, $4 \
    JNE 3(PC) \
    MOVWQZX 0(R11), target \ // unsigned 16 bit
    JMP 6(PC) \
    CMPB AX, $5 \
    JNE 3(PC) \
    MOVBQSX 0(R11), target \ // signed 8 bit
    JMP 2(PC) \
    MOVBQZX 0(R11), target // unsigned 8 bit

#define LOADXMMREG(off, target) \
    MOVLQSX off(R12), AX \
    TESTQ AX, AX \
    JS prepared \
    MOVWQZX AX, R11 \
    SHRL $16, AX \
    ADDQ R13, R11 \
    CMPB AX, $7 \
    JNE 3(PC) \
    MOVSD 0(R11), target \ // float 64bit
    JMP 2(PC) \
    MOVSS 0(R11), target \ // float 32bit


// func asmcall()
TEXT ·asmcall(SB),NOSPLIT,$16 // the 16 fixes the stack alignment which was broken by the call to asmcall
    MOVQ DI, R12      // FRAME (preserved)
    MOVQ 8(R12), R13  // base

    // load register arguments
    LOADREG(40, DI)
    LOADREG(44, SI)
    LOADREG(48, DX)
    LOADREG(52, CX)
    LOADREG(56, R8)
    LOADREG(60, R9)

xmm:

    LOADXMMREG(64, X0)
    LOADXMMREG(68, X1)
    LOADXMMREG(72, X2)
    LOADXMMREG(76, X3)
    LOADXMMREG(80, X4)
    LOADXMMREG(84, X5)
    LOADXMMREG(88, X6)
    LOADXMMREG(92, X7)

prepared:
    // load number of vector registers
    MOVBQZX 112(R12), AX

    // do the actuall call
    CALL (R12)

    // store ret0
    MOVLQSX 96(R12), BX
    TESTQ BX, BX
    JS xmmret0
    MOVWQZX BX, R11
    SHRL $16, BX
    ADDQ R13, R11
    CMPB BX, $0
    JNE 3(PC)
    MOVQ AX, (R11)
    JMP ret1
    CMPB BX, $2
    JGT 3(PC)
    MOVL AX, (R11)
    JMP ret1
    CMPB BX, $4
    JGT 3(PC)
    MOVW AX, (R11)
    JMP ret1
    MOVB AX, (R11)

ret1:

    MOVLQSX 100(R12), BX
    TESTQ BX, BX
    JS DONE
    MOVWQZX BX, R11
    SHRL $16, BX
    ADDQ R13, R11
    CMPB BX, $0
    JNE 3(PC)
    MOVQ DX, (R11)
    JMP ret1
    CMPB BX, $2
    JGT 3(PC)
    MOVL DX, (R11)
    JMP ret1
    CMPB BX, $4
    JGT 3(PC)
    MOVW DX, (R11)
    JMP ret1
    MOVB DX, (R11)

xmmret0:

    MOVLQSX 104(R12), BX
    TESTQ BX, BX
    JS DONE
    MOVWQZX BX, R11
    SHRL $16, BX
    ADDQ R13, R11
    CMPB BX, $7
    JNE 3(PC)
    MOVSD X0, (R11)
    JMP xmmret1
    MOVSS X0, (R11)

xmmret1:

    MOVLQSX 108(R12), BX
    TESTQ BX, BX
    JS DONE
    MOVWQZX BX, R11
    SHRL $16, BX
    ADDQ R13, R11
    CMPB BX, $7
    JNE 3(PC)
    MOVSD X1, (R11)
    JMP xmmret1
    MOVSS X1, (R11)

DONE:
    RET


GLOBL pthread_attr_init__dynload(SB), NOPTR, $8
GLOBL pthread_attr_getstacksize__dynload(SB), NOPTR, $8
GLOBL pthread_attr_destroy__dynload(SB), NOPTR, $8

// runtime has #include "go_asm.h"
// we need to fake those two defines here:
#define g_stack 0
#define stack_lo 0

TEXT x_cgo_init(SB),NOSPLIT,$512 // size_t size (8 byte) + unknown pthread_attr_t - hopefully this is big enough
    MOVQ DI, R12 // g

    // pthread_attr_init(8(SP))
    LEAQ 8(SP), DI
    MOVQ $pthread_attr_init__dynload(SB), R11
    CALL (R11)

    // pthread_attr_init(8(SP), 0(SP))
    LEAQ 8(SP), DI
    LEAQ 0(SP), SI
    MOVQ $pthread_attr_getstacksize__dynload(SB), R11
    CALL (R11)

    // g->stacklo = &size - size + 4096
    LEAQ 0x1000(SP), AX
    SUBQ 0(SP), AX
    MOVQ AX, (g_stack+stack_lo)(R12)

    // pthread_attr_init(8(SP))
    LEAQ 8(SP), DI
    MOVQ $pthread_attr_destroy__dynload(SB), R11
    CALL (R11)

    RET
