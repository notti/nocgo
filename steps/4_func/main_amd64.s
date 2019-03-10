#include "textflag.h"
#include "funcdata.h"
#include "go_asm.h"

TEXT ·callWrapper(SB),NOSPLIT|WRAPPER,$24
    NO_LOCAL_POINTERS
    MOVQ DX, 0(SP)
    LEAQ argframe+0(FP), CX
    MOVQ CX, 8(SP)
    CALL ·fake(SB)
    MOVQ SP, AX
    ADDQ $24+8+8, AX
    MOVQ 0(SP), CX
    MOVQ funcStorage_argsize(CX), CX
    ADDQ CX, AX
    MOVQ 16(SP), BX
    MOVQ BX, (AX)
    RET
