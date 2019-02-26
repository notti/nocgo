#include "textflag.h"
#include "go_asm.h"

// runtime has #include "go_asm.h"
// we need to fake the defines here:
#define slice_array 0
#define slice_len 4
#define slice_cap 8

// func asmcall(spec)
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

    CMPB BX, $const_type64
    JNE 5(PC)
    MOVL Spec_base(SI), BX
    ADDL BX, CX
    MOVL AX, (CX)
    ADDL $4, CX
    MOVL DX, (CX)
    JMP done

    INT $3

done:

    RET
