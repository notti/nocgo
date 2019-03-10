#include "textflag.h"
#include "funcdata.h"
#include "go_asm.h"

TEXT ·callWrapper(SB),NOSPLIT|WRAPPER,$16
    NO_LOCAL_POINTERS
    MOVL DX, 0(SP)
    LEAL argframe+0(FP), CX
    MOVL CX, 4(SP)
    CALL ·fake(SB)
    
    MOVL SP, AX
    ADDL $16+4, AX
    MOVL 0(SP), CX
    MOVL funcStorage_argsize(CX), DX
    ADDL DX, AX
    MOVL funcStorage_retsize(CX), DX
    MOVL SP, CX
    ADDL $8, CX

next:
    MOVL (CX), BX
    MOVL BX, (AX)

    ADDL $4, CX
    ADDL $4, AX
    SUBL $4, DX

    TESTL DX, DX
    JNE next

    RET
