#include "textflag.h"
#include "go_asm.h"
#include "funcdata.h"

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
    JMP 30(PC) \
    CMPB AX, $const_typeS32 \
    JNE 3(PC) \
    MOVLQSX 0(R11), target \ // signed 32 bit
    JMP 26(PC) \
    CMPB AX, $const_typeU32 \
    JNE 3(PC) \
    MOVLQZX 0(R11), target \ // unsigned 32 bit
    JMP 23(PC) \
    CMPB AX, $const_typeS16 \
    JNE 3(PC) \
    MOVWQSX 0(R11), target \ // signed 16 bit
    JMP 20(PC) \
    CMPB AX, $const_typeU16 \
    JNE 3(PC) \
    MOVWQZX 0(R11), target \ // unsigned 16 bit
    JMP 16(PC) \
    CMPB AX, $const_typeS8 \
    JNE 3(PC) \
    MOVBQSX 0(R11), target \ // signed 8 bit
    JMP 12(PC) \
    CMPB AX, $const_typeU8 \
    JNE 3(PC) \
    MOVBQZX 0(R11), target \ // unsigned 8 bit
    JMP 8(PC) \
    CMPB AX, $const_typeCallback \ // callback
    JNE 5(PC) \
    SUBQ R13, R11 \
    MOVQ $callbacks<>(SB), AX \
    MOVQ (AX)(R11*8), target \
    JMP 2(PC) \
    INT $3

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

//func cgocallback(fn, frame unsafe.Pointer, framesize, ctxt uintptr)
TEXT cgocallback(SB),NOSPLIT,$0
    JMP runtime·cgocallback(SB)

// 18*8
// 0x0: fn
// 0x8: frame
// 0x10: framesize
// 0x18: ctx
// 0x20: bp <- callbackArgs
// 0x28: DI
// 0x30: SI
// 0x38: DX <- ret
// 0x40: CX
// 0x48: R8
// 0x50: R9
// 0x58: X0 <- ret
// 0x60: X1
// 0x68: X2
// 0x70: X3
// 0x78: X4
// 0x80: X5
// 0x88: X6
// 0x90: X7
// 0x98: AX <- ret
// 0xA0: which
// 0xA8: BX <- safe
// 0xB0: R12
// 0xB8: R13
// 0xC0: R14
// 0xC8: R15

// need to save BP?
#define CALLBACK(name, id) \
TEXT name(SB),NOSPLIT,$0xD8 \
    MOVQ DI, 0x28(SP) \
    MOVQ SI, 0x30(SP) \
    MOVQ DX, 0x38(SP) \
    MOVQ CX, 0x40(SP) \
    MOVQ R8, 0x48(SP) \
    MOVQ R9, 0x50(SP) \
    TESTB AX, AX \
    JZ skip \
    MOVSD X0, 0x58(SP) \
    MOVSD X1, 0x60(SP) \
    MOVSD X2, 0x68(SP) \
    MOVSD X3, 0x70(SP) \
    MOVSD X4, 0x78(SP) \
    MOVSD X5, 0x80(SP) \
    MOVSD X6, 0x88(SP) \
    MOVSD X7, 0x90(SP) \
skip: \
    MOVQ $id, 0xA0(SP) \
    MOVQ BX, 0xA8(SP) \
    MOVQ R12, 0xB0(SP) \
    MOVQ R13, 0xB8(SP) \
    MOVQ R14, 0xC0(SP) \
    MOVQ R15, 0xC8(SP) \
    LEAQ ·testCallback(SB), AX \
    MOVQ AX, 0(SP) \
    LEAQ 0x20(SP), AX \
    MOVQ AX, 0xD0(SP) \
    LEAQ 0xD0(SP), AX \
    MOVQ AX, 0x8(SP) \
    MOVQ $8, 0x10(SP) \
    MOVQ $0, 0x18(SP) \
    LEAQ arg+0(FP), AX \
    MOVQ AX, 0x20(SP) \
    CALL cgocallback(SB) \
    MOVQ 0x38(SP), DX \
    MOVQ 0x98(SP), AX \
    MOVSD 0x58(SP), X0 \
    MOVQ 0xA8(SP), BX \
    MOVQ 0xB0(SP), R12 \
    MOVQ 0xB8(SP), R13 \
    MOVQ 0xC0(SP), R14 \
    MOVQ 0xC8(SP), R15 \
    RET

CALLBACK(callback0, 0)
CALLBACK(callback1, 1)
CALLBACK(callback2, 2)
CALLBACK(callback3, 3)
CALLBACK(callback4, 4)
CALLBACK(callback5, 5)

DATA callbacks<>+0x00(SB)/8, $callback0(SB)
DATA callbacks<>+0x08(SB)/8, $callback1(SB)
DATA callbacks<>+0x10(SB)/8, $callback2(SB)
DATA callbacks<>+0x18(SB)/8, $callback3(SB)
DATA callbacks<>+0x20(SB)/8, $callback4(SB)
DATA callbacks<>+0x28(SB)/8, $callback5(SB)
GLOBL callbacks<>(SB),RODATA,$48

// pass struct { &args, &spec } to cgocall
TEXT ·callWrapper(SB),NOSPLIT|WRAPPER,$32
    NO_LOCAL_POINTERS
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
    CALL _cgo_topofstack(SB)
    SUBQ AX, R13

    // load number of vector registers
    MOVBQSX spec_rax(R12), AX

    // do the actuall call
    CALL spec_fn(R12)

    MOVQ R14, SP

    MOVQ AX, BX

    // readjust our arguments in case a stack split happened
    CALL _cgo_topofstack(SB) //clobbers AX, CX
    ADDQ AX, R13

    MOVQ BX, AX

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
