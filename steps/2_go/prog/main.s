#include "textflag.h"

// func asmcall3()
TEXT ·asmcall3(SB),NOSPLIT,$8
    MOVQ DI, 0(SP) // for returning the argument
    
    MOVQ 0(DI), R8
    MOVQ 16(DI), SI
    MOVQ 24(DI), DX
    MOVQ 8(DI), DI

    XORQ AX, AX // no floating point

    CALL R8

    MOVQ 0(SP), DI
    MOVQ AX, 0(DI) // return argument in fn - r1 would get optimized away :()

    RET
