#include "textflag.h"
#include "go_asm.h"

// runtime has #include "go_asm.h"
// we need to fake the defines here:
#define slice_array 0
#define slice_len 8
#define slice_cap 16


#define LOADREG(off, target) \
    MOVLQSX spec_intargs+argument__size*off(R12), AX \
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
    MOVLQSX spec_xmmargs+argument__size*off(R12), AX \
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

TEXT ·cgocall(SB),NOSPLIT,$0
    JMP runtime·cgocall(SB)

// pass struct { &args, &spec } to cgocall
TEXT ·callWrapper(SB),NOSPLIT|WRAPPER,$32
    MOVQ DX, 24(SP)
    LEAQ argframe+0(FP), AX
    MOVQ AX, 16(SP)
    LEAQ 16(SP), AX
    MOVQ AX, 8(SP)
    LEAQ asmcall(SB), AX
    MOVQ AX, 0(SP)
    CALL ·cgocall(SB)
    RET

TEXT asmcall(SB),NOSPLIT,$0
    MOVQ 8(DI), R12      // spec (preserved)
    MOVQ 0(DI), R13      // base of args (preserved)
    MOVQ SP, R14         // stack for restoring later on (preserved)

    ANDQ $~0x1F, SP // 32 byte alignment for cdecl (in case someone wants to pass __m256 on the stack)
    // for no __m256 16 byte would be ok
    // this is actually already done by cgocall - but asmcall was called from there and destroys that :(

    MOVQ spec_stack+slice_len(R12), AX // length of stack registers
    TESTQ AX, AX
    JZ reg

    // ok we have stack arguments so let's do that first

    // Fix alignment depending on number of arguments
    MOVQ AX, BX
    ANDQ $3, BX
    SHLQ $3, BX
    SUBQ BX, SP

    MOVQ spec_stack+slice_array(R12), BX

next:
    DECQ AX
    MOVQ (BX)(AX*argument__size), CX
    //check type and push to stack
    MOVWQZX CX, R11
    SHRL $16, CX
    ADDQ R13, R11

#define LOADSTACK(type, instr) \
    CMPB CX, type \
    JNE 7(PC) \
    SUBQ $8, SP \
    instr 0(R11), CX \
    instr CX, 0(SP) \
    TESTQ AX, AX \
    JZ reg \
    JMP next

    LOADSTACK($const_type64, MOVQ)
    LOADSTACK($const_typeU32, MOVL)
    LOADSTACK($const_typeU16, MOVW)
    LOADSTACK($const_typeU8, MOVB)

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
    MOVBQSX spec_rax(R12), AX

    // do the actuall call
    CALL spec_fn(R12)

    MOVQ R14, SP

    // TODO: check R13, if it still points to the correct stack! (could happen if we have a callback into go that splits the stack)

    // store ret
    MOVLQSX spec_ret(R12), BX
    TESTQ BX, BX
    JS DONE
    MOVWQZX BX, R11
    SHRL $16, BX
    ADDQ R13, R11

    CMPB BX, $const_type64
    JNE 3(PC)
    MOVQ AX, (R11)
    JMP DONE

    CMPB BX, $const_typeU32
    JNE 3(PC)
    MOVL AX, (R11)
    JMP DONE

    CMPB BX, $const_typeU16
    JNE 3(PC)
    MOVW AX, (R11)
    JMP DONE

    CMPB BX, $const_typeU8
    JNE 3(PC)
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
