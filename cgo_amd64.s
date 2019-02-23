/*
    This contains the code needed to provide a C-like environment inside go
    (e.g., pthread setup, so TLS is done correctly etc.)
    In go this is done by cgo (runtime/cgo) with C-code.
    We can't do that here -> we'll use go assembly.

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

    Our special linker will set iscgo to 1 and set the _cgo_* stuff to x_cgo_*

    We didn't implement the following functions, since it's only needed for TSAN:
     * _cgo_mmap
     * _cgo_munmap
     * _cgo_sigaction
*/
#include "textflag.h"
// runtime has #include "go_asm.h"
// we need to fake the defines here:
#define g_stack 0
#define stack_lo 0
#define stack_hi 8

#define timespec__size 16
#define sigset_t__size 128
#define pthread_t__size 8
#define pthread_attr_t__size 56
#define size_t__size 8

// symbols we import from C
GLOBL pthread_attr_init__dynload(SB), NOPTR, $8
GLOBL pthread_attr_getstacksize__dynload(SB), NOPTR, $8
GLOBL pthread_attr_destroy__dynload(SB), NOPTR, $8
GLOBL pthread_sigmask__dynload(SB), NOPTR, $8
GLOBL pthread_create__dynload(SB), NOPTR, $8
GLOBL pthread_detach__dynload(SB), NOPTR, $8
GLOBL setenv__dynload(SB), NOPTR, $8
GLOBL unsetenv__dynload(SB), NOPTR, $8
GLOBL malloc__dynload(SB), NOPTR, $8
GLOBL free__dynload(SB), NOPTR, $8
GLOBL nanosleep__dynload(SB), NOPTR, $8
GLOBL sigfillset__dynload(SB), NOPTR, $8

// storage for setg
GLOBL setg_gcc(SB), NOPTR, $8

// _cgo_init(G *g, void (*setg)(void*)) (runtime/cgo/gcc_linux_amd64.c)
// This get's called during startup, adjusts stacklo, and provides a pointer to setg_gcc for us
// Additionally, if we set _cgo_init to non-null, go won't do it's own TLS setup
//
// 0(SP): size_t size (8 byte)
// 8(SP): pthread_attr_t attr (56)
TEXT x_cgo_init(SB),NOSPLIT,$64
    MOVQ DI, R12           // g
    MOVQ SI, setg_gcc(SB)  // setg

    // pthread_attr_init(&attr = &8(SP))
    LEAQ 8(SP), DI
    MOVQ $pthread_attr_init__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // pthread_attr_getstacksize(&size = &0(SP), &attr = &8(SP))
    LEAQ 8(SP), DI
    LEAQ 0(SP), SI
    MOVQ $pthread_attr_getstacksize__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // g->stacklo = &size - size + 4096
    LEAQ 0x1000(SP), AX
    SUBQ 0(SP), AX
    MOVQ AX, stack_lo(R12)

    // pthread_attr_init(&attr = &8(SP))
    LEAQ 8(SP), DI
    MOVQ $pthread_attr_destroy__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    RET

// _cgo_thread_start is split into three parts in cgo since only one part
// is system dependent - we keep the split here, since it is easier to code :)

// _cgo_thread_start(ThreadStart *arg) (runtime/cgo/gcc_util.c)
// This get's called instead of the go code for creating new threads
// -> pthread_* stuff is used, so threads are setup correctly for C
// If this is missing, TLS is only setup correctly on thread 1
//
// ThreadStart is size 24 
TEXT x_cgo_thread_start(SB),NOSPLIT,$0
    MOVQ DI, R12

    // AX = malloc(24)
    MOVQ $24, DI
    MOVQ $malloc__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    TESTQ AX, AX
    JNZ OK
    INT $3 // malloc failed - final code should write some error message

OK:
    // *ax = *ts
    MOVQ 0(R12), DI
    MOVQ DI, 0(AX)
    MOVQ 8(R12), DI
    MOVQ DI, 8(AX)
    MOVQ 16(R12), DI
    MOVQ DI, 16(AX)

    // _cgo_sys_thread_start(ax)
    MOVQ AX, DI
    XORQ AX, AX
    CALL _cgo_sys_thread_start(SB)
    RET

// _cgo_sys_thread_start(ThreadStart *arg) (runtime/cgo/gcc_linux_amd64.c)
//   0(SP) size_t size (size 8)
//   8(SP) sigset_t ign (128)
// 136(SP) sigset_t oset (128)
// 264(SP) attr_t attr (56)
TEXT _cgo_sys_thread_start(SB),NOSPLIT,$320
    MOVQ DI, R12

    // sigfillset(&ign = &8(SP))
    LEAQ 8(SP), DI
    MOVQ $sigfillset__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // pthread_sigmask(SIG_SETMASK = 2, &ign = &8(SP), &oset = &136(SP))
    MOVQ $2, DI
    LEAQ 8(SP), SI
    LEAQ 136(SP), DX
    MOVQ $pthread_sigmask__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // pthread_attr_init(&attr = &264(SP))
    LEAQ 264(SP), DI
    MOVQ $pthread_attr_init__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // pthread_attr_getstacksize(&attr = &264(SP), &size = &0(SP))
    LEAQ 264(SP), DI
    LEAQ 0(SP), SI
    MOVQ $pthread_attr_getstacksize__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // ts->g->stack_high = size
    MOVQ 0(R12), AX
    MOVQ 0(SP), BX
    MOVQ BX, stack_hi(AX) //FIXME: check

    // In the cgocode this has *pthread_t as first argument, but we don't use this in this function...
    // And by not doing that it is simpler here
    // R12 = _cgo_try_pthread_create(&attr = 264(SP), threadentry, &ts)
    LEAQ 264(SP), DI
    MOVQ $threadentry(SB), SI
    MOVQ R12, DX
    XORQ AX, AX
    // We didn't clobber any registers here -> no saving required
    CALL _cgo_try_pthread_create(SB)

    MOVQ AX, R12

    // pthread_sigmask(SIG_SETMASK = 2, &oset = &136(SP), nil)
    MOVQ $2, DI
    LEAQ 136(SP), SI
    XORQ DX, DX
    MOVQ $pthread_sigmask__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // if R12 != 0 -> fail
    TESTQ R12, R12
    JZ OK

    INT $3 // pthread_create failed

OK:
    RET

// _cgo_try_pthread_create(*attr, *threadentry, *ts) (runtime/cgo/gcc_libinit.c)
// 0(SP) pthread_t (8)
// 8(SP) timespec_t (16)
TEXT _cgo_try_pthread_create(SB),NOSPLIT,$24
    MOVQ DI, R12
    MOVQ SI, R13
    MOVQ DX, R14
    
    MOVQ $0, BX // tries

AGAIN:
    // AX = pthread_create(&pthread_t = &0(SP), attr, threadentry, ts)
    LEAQ 0(SP), DI // &pthread_t
    MOVQ R12, SI
    MOVQ R13, DX
    MOVQ R14, CX
    MOVQ $pthread_create__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    TESTQ AX, AX
    JNZ FAILED

    // if AX == 0
    // pthread_detach(pthread_t)
    MOVQ 0(SP), DI
    MOVQ $pthread_detach__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // return 0
    XORQ AX, AX
    RET

    // AX != 0:
FAILED:
    CMPQ AX, $11 //EAGAIN
    JEQ WAIT

    // AX != EAGAIN -> return the error
    RET

WAIT:
    // AX == EAGAIN

    // ts.tv_sec = 0
    MOVQ $0, 8(SP)
    // ts.tv_nsec = (1+tries)*1000*1000
    MOVQ BX, AX
    INCQ AX
    IMULQ $(1000*1000), AX
    MOVQ AX, 16(SP)
    
    // nanosleep(&ts = &8(SP), nil)
    LEAQ 8(SP), SI
    XORQ DI, DI
    MOVQ $nanosleep__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // tries++
    INCQ BX
    CMPQ BX, $20
    JLT AGAIN

    // if tries >= 20 return eagain

    MOVQ $11, AX //EAGAIN
    RET

// threadentry(ThreadStart *ts) (runtime/cgo/gcc_linux_amd64.c)
// 0(SP) ThreadStart (size 24)
TEXT threadentry(SB),NOSPLIT,$24 //ThreadStart
    // 0(SP) = *ts
    MOVQ 0(DI), AX
    MOVQ AX, 0(SP)
    MOVQ 8(DI), AX
    MOVQ AX, 8(SP)
    MOVQ 16(DI), AX
    MOVQ AX, 16(SP)

    // free(ts)
    // ts is still the first argument
    MOVQ $free__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)

    // setg_gcc(ts->g)
    MOVQ 0(SP), DI
    MOVQ $setg_gcc(SB), R11
    XORQ AX, AX
    CALL (R11)

    // ts.fn
    // in gcc_linux_amd64.c this uses crosscall, which we don't need here (we don't care about any register after this point)
    // also no need for clearing AX -> we call into go
    MOVQ 16(SP), R11
    CALL R11

    // return 0
    XORQ AX, AX
    RET

// _cgo_notify_runtime_init_done (runtime/cgo/gcc_libinit.c)
// do nothing - we don't support being a library
TEXT x_cgo_notify_runtime_init_done(SB),NOSPLIT,$0
    RET

// _cgo_setenv(char **arg) (runtime/cgo/gcc_setenv.c)
TEXT x_cgo_setenv(SB),NOSPLIT,$0
    // setenv(arg[0], arg[1], 1)
    MOVQ $1, DX
    MOVQ 8(DI), SI
    MOVQ (DI), DI
    MOVQ $setenv__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)
    RET

// _cgo_unsetenv(char *arg) (runtime/cgo/gcc_setenv.c)
TEXT x_cgo_unsetenv(SB),NOSPLIT,$0
    // unsetenv(arg)
    MOVQ $setenv__dynload(SB), R11
    XORQ AX, AX
    CALL (R11)
    RET

// _cgo_callers(uintptr_t sig, void *info, void *context, void (*cgoTraceback)(struct cgoTracebackArg*), uintptr_t* cgoCallers, void (*sigtramp)(uintptr_t, void*, void*)) (runtime/cgo/gcc_traceback.c)
TEXT x_cgo_callers(SB),NOSPLIT,$32
    PUSHQ R9
    PUSHQ DI
    PUSHQ SI
    PUSHQ DX
    
    // arg.Context = 0
    MOVQ $0, (SP)
    // arg.SigContext = context
    MOVQ DX, 8(SP)
    // arg.Buf = cgoCallers
    MOVQ R8, 16(SP)
    // arg.Max = 32
    MOVQ $32, 24(SP)

    // cgoTraceback(&arg)
    MOVQ SP, DI
    XORQ AX, AX
    CALL (CX)

    // sigtramp(sig, info, context)
    POPQ DX
    POPQ SI
    POPQ DI
    POPQ R9
    XORQ AX, AX
    CALL (R9)
    RET
