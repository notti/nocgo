/*
    trampoline for emulating required C functions for cgo in go (see cgoGlue.go)
    (we convert cdecl calling convention to go and vice-versa)

    Since we're called from go and call into C we can cheat a bit with the calling conventions:
     - in go all the registers are caller saved
     - in C we have a couple of callee saved registers

    => we can use BX, R12, R13, R14, R15 instead of the stack

    C Calling convention cdecl used here (we only need integer args):
    1. arg: DI
    2. arg: SI
    3. arg: DX
    4. arg: CX
    5. arg: R8
    6. arg: R9
    We don't need floats with these functions -> AX=0
    return value will be in AX
*/
#include "textflag.h"

TEXT x_cgo_init_trampoline(SB),NOSPLIT,$16
    MOVQ DI, 0(SP)
    MOVQ SI, 8(SP)
    CALL ·x_cgo_init(SB)
    RET

TEXT x_cgo_thread_start_trampoline(SB),NOSPLIT,$8
    MOVQ DI, 0(SP)
    CALL ·x_cgo_thread_start(SB)
    RET

TEXT ·threadentry_trampoline(SB),NOSPLIT,$16
    MOVQ DI, 0(SP)
    CALL ·threadentry(SB)
    MOVQ 8(SP), AX
    RET

// func setg_trampoline(setg uintptr, g uintptr)
TEXT ·setg_trampoline(SB),NOSPLIT,$0-16
    MOVQ g+8(FP), DI
    MOVQ setg+0(FP), AX
    CALL AX
    RET

// _cgo_notify_runtime_init_done (runtime/cgo/gcc_libinit.c)
TEXT x_cgo_notify_runtime_init_done_trampoline(SB),NOSPLIT,$0
    CALL ·x_cgo_notify_runtime_init_done(SB)
    RET

// _cgo_setenv(char **arg) (runtime/cgo/gcc_setenv.c)
TEXT x_cgo_setenv_trampoline(SB),NOSPLIT,$8
    MOVQ DI, 0(SP)
    CALL ·x_cgo_setenv(SB)
    RET

// _cgo_unsetenv(char *arg) (runtime/cgo/gcc_setenv.c)
TEXT x_cgo_unsetenv_trampoline(SB),NOSPLIT,$8
    MOVQ DI, 0(SP)
    CALL ·x_cgo_unsetenv(SB)
    RET

// func asmlibccall6(fn, n, args uintptr) uintptr
TEXT ·asmlibccall6(SB),NOSPLIT,$0-32
    MOVQ n+8(FP), AX

    TESTQ AX, AX
    JZ skipargs

    MOVQ args+16(FP), AX
    MOVQ 0(AX), DI
    MOVQ 8(AX), SI
    MOVQ 16(AX), DX
    MOVQ 24(AX), CX
    MOVQ 32(AX), R8
    MOVQ 40(AX), R9

skipargs:

    MOVQ fn+0(FP), BX

    MOVQ SP, R12
    ANDQ $~0xF, SP // 16 byte alignment for cdecl

    XORQ AX, AX // no fp arguments
    CALL BX

    MOVQ R12, SP

    MOVQ AX, ret+24(FP)

    RET
