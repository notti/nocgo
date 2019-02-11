#include "textflag.h"

/*
    Frame layout:
    fn
    base
    regarg0
    regarg1
    regarg2
    regarg3
    regarg4
    regarg5
    rax
    ret0
    ret1
*/

#define LOADREG(off, target) \
    MOVQ off(R12), R11 \
    ADDQ R13, R11 \
    MOVQ 0(R11), target

// func asmcall3()
TEXT ·asmcall3(SB),NOSPLIT,$0
    MOVQ DI, R12     // FRAME (preserved)
    MOVQ 8(R12), R13  // base
    
    // load register arguments
    LOADREG(16, DI)
    LOADREG(24, SI)
    LOADREG(32, DX)
    LOADREG(40, CX)
    LOADREG(48, R8)
    LOADREG(56, R9)

    // load number of vector registers
    MOVQ 64(R12), AX

    // do the actuall call
    CALL (R12)

    // store ret0
    MOVQ 72(R12), BX
    TESTQ BX, BX
    JZ DONE
    ADDQ R13, BX
    MOVQ AX, (BX)

    // store ret1
    MOVQ 80(R12), BX
    TESTQ BX, BX
    JZ DONE
    ADDQ R13, BX
    MOVQ DX, (BX)

DONE:
    RET
