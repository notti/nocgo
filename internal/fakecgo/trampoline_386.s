/*
    See trampoline_amd64.s for explanations.

    => we can use SI, DI, BX instead of the stack

    C Calling convention cdecl used here (we only need integer args):
    Pass values on the stack in right to left order
    return value will be in AX

    If there are no return values calling conventions are the same -> NOFRAME + JMP
*/
#include "textflag.h"

TEXT x_cgo_init_trampoline(SB),NOSPLIT|NOFRAME,$0
    JMP ·x_cgo_init(SB)
    RET

TEXT x_cgo_thread_start_trampoline(SB),NOSPLIT|NOFRAME,$0
    JMP ·x_cgo_thread_start(SB)
    RET

TEXT ·threadentry_trampoline(SB),NOSPLIT,$8
    MOVL arg0+0(FP), AX
    MOVL AX, 0(SP)
    CALL ·threadentry(SB)
    MOVL 4(SP), AX
    RET

// func setg_trampoline(setg uintptr, g uintptr)
TEXT ·setg_trampoline(SB),NOSPLIT,$4-8
    MOVL g+4(FP), AX
    MOVL AX, 0(SP)
    MOVL setg+0(FP), AX
    CALL AX
    RET

TEXT x_cgo_notify_runtime_init_done_trampoline(SB),NOSPLIT|NOFRAME,$0
    JMP ·x_cgo_notify_runtime_init_done(SB)
    RET

TEXT x_cgo_setenv_trampoline(SB),NOSPLIT|NOFRAME,$0
    JMP ·x_cgo_setenv(SB)
    RET

TEXT x_cgo_unsetenv_trampoline(SB),NOSPLIT|NOFRAME,$0
    JMP ·x_cgo_unsetenv(SB)
    RET

// func asmlibccall6(fn, n, args uintptr) uintptr
TEXT ·asmlibccall6(SB),NOSPLIT,$0-16
    MOVL SP, DI
    MOVL n+4(FP), AX
    MOVL args+8(FP), BX
    MOVL fn+0(FP), DX

    TESTL AX, AX
    JZ finished

next:
    DECL AX
    MOVL (BX)(AX*4), CX
    SUBL $4, SP
    MOVL CX, 0(SP)

    TESTL AX, AX
    JZ finished
    JMP next

finished:
    CALL DX

    MOVL DI, SP

    MOVL AX, ret+12(FP)

    RET
