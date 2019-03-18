#include "textflag.h"
#include "go_asm.h"

// runtime has #include "go_asm.h"
// we need to fake the defines here:
#define slice_array 0
#define slice_len 4
#define slice_cap 8

TEXT ·cgocall(SB),NOSPLIT,$0
    JMP runtime·cgocall(SB)

// pass struct { &args, &spec } to cgocall
TEXT ·callWrapper(SB),NOSPLIT|WRAPPER,$16
    MOVL DX, 12(SP)
    LEAL argframe+0(FP), AX
    MOVL AX, 8(SP)
    LEAL 8(SP), AX
    MOVL AX, 4(SP)
    LEAL asmcall(SB), AX
    MOVL AX, 0(SP)
    CALL ·cgocall(SB)
    RET

TEXT asmcall(SB),NOSPLIT,$0-4
    MOVL frame+0(FP), SI // &args, &spec (preserved)
    MOVL SP, DI          // STACK (preserved)

    MOVL 4(SI), AX       // spec
    MOVL spec_stack+slice_len(AX), BX
    TESTL BX, BX
    JZ prepared

    MOVL spec_stack+slice_array(AX), DX

next:
    DECL BX
    MOVL (DX)(BX*argument__size), AX
    MOVWLZX AX, CX
    SHRL $16, AX

#define TYPE(which, instr) \
    CMPB AX, which \
    JNE 9(PC) \
    SUBL $4, SP \
    MOVL 0(SI), AX \
    ADDL CX, AX \
    instr 0(AX), AX \
    MOVL AX, 0(SP) \
    TESTL BX, BX \
    JZ prepared \
    JMP next

    TYPE($const_type32, MOVL)
    TYPE($const_type16, MOVW)
    TYPE($const_type8, MOVB)
    INT $3

prepared:

    MOVL 4(SI), AX
    CALL spec_fn(AX)

    // return value in AX, DX, F0 <- DON'T USE THESE

    MOVL DI, SP // restore stack

    MOVL 4(SI), DI // DI: spec

    MOVL spec_ret(DI), BX
    TESTL BX, BX
    JS done

    MOVL 0(SI), SI // SI: args

    // TODO: check SI, if it still points to the correct stack! (could happen if we have a callback into go that splits the stack)

    MOVWLZX BX, CX
    SHRL $16, BX

    CMPB BX, $const_type32
    JNE 4(PC)
    ADDL SI, CX
    MOVL AX, (CX)
    JMP done

    CMPB BX, $const_type16
    JNE 4(PC)
    ADDL SI, CX
    MOVW AX, (CX)
    JMP done

    CMPB BX, $const_type8
    JNE 4(PC)
    ADDL SI, CX
    MOVB AX, (CX)
    JMP done

    CMPB BX, $const_typeFloat
    JNE 4(PC)
    ADDL SI, CX
    FMOVF F0, (CX)
    JMP done

    CMPB BX, $const_typeDouble
    JNE 4(PC)
    ADDL SI, CX
    FMOVD F0, (CX)
    JMP done

    CMPB BX, $const_type64
    JNE 6(PC)
    ADDL SI, CX
    MOVL AX, (CX)
    ADDL $4, CX
    MOVL DX, (CX)
    JMP done

    INT $3

done:

    RET
