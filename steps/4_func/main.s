#include "textflag.h"
#include "funcdata.h"

TEXT ·callWrapper(SB),NOSPLIT|WRAPPER,$24
    NO_LOCAL_POINTERS
    MOVQ DX, 0(SP)
    LEAQ argframe+0(FP), CX
    MOVQ CX, 8(SP)
    CALL ·fake(SB)
    MOVQ 16(SP), AX
    MOVQ AX, ret+8(FP)
    RET
