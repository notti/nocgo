#include "textflag.h"
#include "go_asm.h"

// runtime has #include "go_asm.h"
// we need to fake the defines here:
#define g_stack 0
#define stack_lo 0
#define slice_array 0
#define slice_len 4
#define slice_cap 8

#define TYPE(which, instr) \
    CMPB AX, which \
    JNE 9(PC) \
    SUBL $4, SP \
    MOVL Spec_base(SI), AX \
    ADDL CX, AX \
    instr 0(AX), AX \
    MOVL AX, 0(SP) \
    TESTL BX, BX \
    JZ prepared \
    JMP next

// func asmcall()
TEXT ·asmcall(SB),NOSPLIT,$0-4
    MOVL frame+0(FP), SI // FRAME (preserved)
    MOVL SP, DI          // STACK (preserved)

    MOVL Spec_stack+slice_len(SI), BX
    TESTL BX, BX
    JZ prepared

    MOVL Spec_stack+slice_array(SI), DX

next:
    DECL BX
    MOVL (DX)(BX*argument__size), AX
    MOVWLZX AX, CX
    SHRL $16, AX

    TYPE($const_type32, MOVL)
    TYPE($const_typeS16, MOVWLSX)
    TYPE($const_typeU16, MOVWLZX)
    TYPE($const_typeS8, MOVBLSX)
    TYPE($const_typeU8, MOVBLZX)
    INT $3

prepared:

    CALL (SI)

    MOVL DI, SP
    
    MOVL Spec_ret(SI), BX
    TESTL BX, BX
    JS done
    MOVWLZX BX, CX
    SHRL $16, BX
    CMPB BX, $0
    JNE 5(PC)
    MOVL Spec_base(SI), BX
    ADDL BX, CX
    MOVL AX, (CX)
    JMP done

    CMPB BX, $2
    JGT 5(PC)
    MOVL Spec_base(SI), BX
    ADDL BX, CX
    MOVW AX, (CX)
    JMP done

    CMPB BX, $4
    JGT 5(PC)
    MOVL Spec_base(SI), BX
    ADDL BX, CX
    MOVB AX, (CX)
    JMP done

    CMPB BX, $const_typeFloat
    JNE 5(PC)
    MOVL Spec_base(SI), BX
    ADDL BX, CX
    FMOVF F0, (CX)
    JMP done

    CMPB BX, $const_typeDouble
    JNE 5(PC)
    MOVL Spec_base(SI), BX
    ADDL BX, CX
    FMOVD F0, (CX)
    JMP done

    INT $3

done:

/*
prepared:
    // load number of vector registers
    MOVBQZX Spec_rax(R12), AX

    // do the actuall call
    CALL (R12)

    MOVQ R14, SP

    // store ret0
    MOVLQSX Spec_ret0(R12), BX
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
    // store ret1
    MOVLQSX Spec_ret1(R12), BX
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
    // store xmmret0
    MOVLQSX Spec_xmmret0(R12), BX
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
    // store xmmret1
    MOVLQSX Spec_xmmret1(R12), BX
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

DONE:*/
    RET


GLOBL pthread_attr_init__dynload(SB), NOPTR, $4
GLOBL pthread_attr_getstacksize__dynload(SB), NOPTR, $4
GLOBL pthread_attr_destroy__dynload(SB), NOPTR, $4

TEXT x_cgo_init(SB),NOSPLIT,$512-4 // size_t size (8 byte) + unknown pthread_attr_t - hopefully this is big enough
    // pthread_attr_init(4(SP))
    LEAL 4(SP), AX
    PUSHL AX
    MOVL $pthread_attr_init__dynload(SB), AX
    CALL (AX)
    POPL AX

    // pthread_attr_init(4(SP), 0(SP))
    LEAL 0(SP), AX
    PUSHL AX
    LEAL 4(SP), AX
    PUSHL AX
    MOVL $pthread_attr_getstacksize__dynload(SB), AX
    CALL (AX)
    POPL AX
    POPL AX

    // g->stacklo = &size - size + 4096
    LEAL 0x1000(SP), AX
    SUBL 0(SP), AX
    MOVL g+0(FP), BX
    MOVL AX, (g_stack+stack_lo)(BX)

    // pthread_attr_destroy(4(SP))
    LEAL 4(SP), AX
    PUSHL AX
    MOVL $pthread_attr_destroy__dynload(SB), AX
    CALL (AX)
    POPL AX

    RET
