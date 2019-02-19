#include "textflag.h"
#include "go_asm.h"

// runtime has #include "go_asm.h"
// we need to fake the defines here:
#define g_stack 0
#define stack_lo 0
#define slice_array 0
#define slice_len 8
#define slice_cap 16


/*
    int:
    type64:     movq              64 bit
    typeS32:    movlqsx    signed 32 bit
    typeU32:    movlqzx  unsigned 32 bit
    typeS16:    movwqsx    signed 16 bit
    typeU16:    movwqzx  unsigned 16 bit
    typeS8:     movbqsx    signed 8  bit
    typeU8:     movbqzx  unsigned 8  bit

    float:
    typeDouble: movsd             64 bit
    typeFloat:  movss             32 bit
*/

#define LOADREG(off, target) \
    MOVLQSX Spec_intargs+argument__size*off(R12), AX \
    TESTQ AX, AX \
    JS xmm \
    MOVWQZX AX, R11 \
    SHRL $16, AX \
    ADDQ R13, R11 \
    CMPB AX, $const_type64 \
    JNE 3(PC) \
    MOVQ 0(R11), target \ // 64bit
    JMP 20(PC) \
    CMPB AX, $const_typeS32 \
    JNE 3(PC) \
    MOVLQSX 0(R11), target \ // signed 32 bit
    JMP 18(PC) \
    CMPB AX, $const_typeU32 \
    JNE 3(PC) \
    MOVLQZX 0(R11), target \ // unsigned 32 bit
    JMP 14(PC) \
    CMPB AX, $const_typeS16 \
    JNE 3(PC) \
    MOVWQSX 0(R11), target \ // signed 16 bit
    JMP 10(PC) \
    CMPB AX, $const_typeU16 \
    JNE 3(PC) \
    MOVWQZX 0(R11), target \ // unsigned 16 bit
    JMP 6(PC) \
    CMPB AX, $const_typeS8 \
    JNE 3(PC) \
    MOVBQSX 0(R11), target \ // signed 8 bit
    JMP 2(PC) \
    MOVBQZX 0(R11), target // unsigned 8 bit

#define LOADXMMREG(off, target) \
    MOVLQSX Spec_xmmargs+argument__size*off(R12), AX \
    TESTQ AX, AX \
    JS prepared \
    MOVWQZX AX, R11 \
    SHRL $16, AX \
    ADDQ R13, R11 \
    CMPB AX, $const_typeDouble \
    JNE 3(PC) \
    MOVSD 0(R11), target \ // float 64bit
    JMP 2(PC) \
    MOVSS 0(R11), target \ // float 32bit


// func asmcall()
TEXT ·asmcall(SB),NOSPLIT,$0
    MOVQ DI, R12      // FRAME (preserved)
    MOVQ Spec_base(R12), R13  // base
    MOVQ SP, R14 // stack

    ANDQ $~0x1F, SP // 32 byte alignment for cdecl (in case someone wants to pass __m256 on the stack)
    // for no __m256 16 byte would be ok
    // this is actually already done by cgocall - but asmcall was called from there and destroys that :(

    MOVQ Spec_stack+slice_len(R12), AX // length of stack registers
    TESTQ AX, AX
    JZ reg

    // ok we have stack arguments so let's do that first

    // Fix alignment depending on number of arguments
    MOVQ AX, BX
    ANDQ $3, BX
    SHLQ $3, BX
    SUBQ BX, SP

    MOVQ Spec_stack+slice_array(R12), BX

next:
    DECQ AX
    MOVQ (BX)(AX*argument__size), CX
    //check type and push to stack
    MOVWQZX CX, R11
    SHRL $16, CX
    ADDQ R13, R11

#define LOADSTACK(type, instr, tmp) \
    CMPB CX, type \
    JNE 7(PC) \
    SUBQ $8, SP \
    instr 0(R11), tmp \
    MOVQ tmp, 0(SP) \
    TESTQ AX, AX \
    JZ reg \
    JMP next

#define LOADSTACKINT(type, instr) LOADSTACK(type, instr, CX)
#define LOADSTACKXMM(type, instr) LOADSTACK(type, instr, X0)

    LOADSTACKINT($const_type64, MOVQ)
    LOADSTACKINT($const_typeS32, MOVLQSX)
    LOADSTACKINT($const_typeU32, MOVLQZX)
    LOADSTACKINT($const_typeS16, MOVWQSX)
    LOADSTACKINT($const_typeU16, MOVWQZX)
    LOADSTACKINT($const_typeS8, MOVBQSX)
    LOADSTACKINT($const_typeU8, MOVBQZX)

    LOADSTACKXMM($const_typeDouble, MOVSD)
    LOADSTACKXMM($const_typeFloat, MOVSS)

    INT $3

reg:
    // load register arguments
    LOADREG(0, DI)
    LOADREG(1, SI)
    LOADREG(2, DX)
    LOADREG(3, CX)
    LOADREG(4, R8)
    LOADREG(5, R9)

xmm:
    // load xmm arguments
    LOADXMMREG(0, X0)
    LOADXMMREG(1, X1)
    LOADXMMREG(2, X2)
    LOADXMMREG(3, X3)
    LOADXMMREG(4, X4)
    LOADXMMREG(5, X5)
    LOADXMMREG(6, X6)
    LOADXMMREG(7, X7)

prepared:
    // load number of vector registers
    MOVBQZX Spec_rax(R12), AX

    // do the actuall call
    CALL (R12)

    MOVQ R14, SP

    // store ret
    MOVLQSX Spec_ret(R12), BX
    TESTQ BX, BX
    JS DONE
    MOVWQZX BX, R11
    SHRL $16, BX
    ADDQ R13, R11

    CMPB BX, $0
    JNE 3(PC)
    MOVQ AX, (R11)
    JMP DONE

    CMPB BX, $2
    JGT 3(PC)
    MOVL AX, (R11)
    JMP DONE

    CMPB BX, $4
    JGT 3(PC)
    MOVW AX, (R11)
    JMP DONE

    CMPB BX, $6
    JGT 3(PC)
    MOVB AX, (R11)
    JMP DONE

    CMPB BX, $const_typeDouble
    JNE 3(PC)
    MOVSD X0, (R11)
    JMP DONE

    CMPB BX, $const_typeFloat
    JNE 3(PC)
    MOVSS X0, (R11)
    JMP DONE

    INT $3

DONE:
    RET


GLOBL pthread_attr_init__dynload(SB), NOPTR, $8
GLOBL pthread_attr_getstacksize__dynload(SB), NOPTR, $8
GLOBL pthread_attr_destroy__dynload(SB), NOPTR, $8

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
