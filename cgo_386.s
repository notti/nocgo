/*
    See cgo_amd64.s for explanations.

    => we can use SI, DI, BX instead of the stack

    C Calling convention cdecl used here (we only need integer args):
    Pass values on the stack in right to left order
    return value will be in AX
*/
#include "textflag.h"
// runtime has #include "go_asm.h"
// we need to fake the defines here:
#define g_stack 0
#define stack_lo 0
#define stack_hi 4

#define timespec__size 8
#define sigset_t__size 128
#define pthread_t__size 4
#define pthread_attr_t__size 36
#define size_t__size 4

// symbols we import from C
GLOBL pthread_attr_init__dynload(SB), NOPTR, $4
GLOBL pthread_attr_getstacksize__dynload(SB), NOPTR, $4
GLOBL pthread_attr_destroy__dynload(SB), NOPTR, $4
GLOBL pthread_sigmask__dynload(SB), NOPTR, $4
GLOBL pthread_create__dynload(SB), NOPTR, $4
GLOBL pthread_detach__dynload(SB), NOPTR, $4
GLOBL setenv__dynload(SB), NOPTR, $4
GLOBL unsetenv__dynload(SB), NOPTR, $4
GLOBL malloc__dynload(SB), NOPTR, $4
GLOBL free__dynload(SB), NOPTR, $4
GLOBL nanosleep__dynload(SB), NOPTR, $4
GLOBL sigfillset__dynload(SB), NOPTR, $4

// storage for setg
GLOBL x_setg_gcc(SB), NOPTR, $4

// _cgo_init(G *g, void (*setg)(void*)) (runtime/cgo/gcc_linux_amd64.c)
//
//  0(SP): arg0
//  4(SP): arg1
//  8(SP): size_t size (4)
// 12(SP): pthread_attr_t attr (36)
TEXT x_cgo_init(SB),NOSPLIT,$48-8
    // setg_gcc = setg
    MOVL setg_gcc+4(FP), AX
    MOVL AX, x_setg_gcc(SB)

    //FIXME: clear attr in stack before calling init (workaround for getting old pthread_attr_init symbol)
    XORL AX, AX
    MOVL AX, 12(SP)
    MOVL AX, 16(SP)
    MOVL AX, 20(SP)
    MOVL AX, 24(SP)
    MOVL AX, 28(SP)
    MOVL AX, 32(SP)
    MOVL AX, 36(SP)
    MOVL AX, 40(SP)
    MOVL AX, 44(SP)

    // pthread_attr_init(&attr = &12(SP))
    LEAL 12(SP), AX
    MOVL AX, 0(SP)
    MOVL $pthread_attr_init__dynload(SB), AX
    CALL (AX)

    // pthread_attr_init(&attr = &12(SP), &size = &8(SP))
    LEAL 12(SP), AX
    MOVL AX, 0(SP)
    LEAL 8(SP), AX
    MOVL AX, 4(SP)
    MOVL $pthread_attr_getstacksize__dynload(SB), AX
    CALL (AX)

    // g->stacklo = &size - size + 4096
    LEAL 0x1000(SP), AX
    SUBL 8(SP), AX
    MOVL g+0(FP), BX
    MOVL AX, (g_stack+stack_lo)(BX)

    // pthread_attr_destroy(4(SP))
    LEAL 4(SP), AX
    MOVL AX, 0(SP)
    MOVL $pthread_attr_destroy__dynload(SB), AX
    CALL (AX)

    RET


// _cgo_thread_start(ThreadStart *arg) (runtime/cgo/gcc_util.c)
//
// ThreadStart is size 12
TEXT x_cgo_thread_start(SB),NOSPLIT,$4-4
    // AX = malloc(12)
    MOVL $12, 0(SP)
    MOVL $malloc__dynload(SB), AX
    CALL (AX)

    TESTL AX, AX
    JNZ OK
    INT $3 // malloc failed - final code should write some error message

OK:
    // *ax = *ts
    MOVL ts+0(FP), BX
    MOVL 0(BX), CX
    MOVL CX, 0(AX)
    MOVL 4(BX), CX
    MOVL CX, 4(AX)
    MOVL 8(BX), CX
    MOVL CX, 8(AX)

    // _cgo_sys_thread_start(ax)
    MOVL AX, 0(SP)
    CALL _cgo_sys_thread_start(SB)
    RET

// _cgo_sys_thread_start(ThreadStart *arg) (runtime/cgo/gcc_linux_amd64.c)
//  0(SP) arg0
//  4(SP) arg1
//  8(SP) arg2
// 12(SP) size_t size (size 4)
// 16(SP) sigset_t ign (128)
//148(SP) sigset_t oset (128)
//276(SP) attr_t attr (36)
TEXT _cgo_sys_thread_start(SB),NOSPLIT,$312-4
    //FIXME: clear attr in stack before calling init (workaround for getting old pthread_attr_init symbol)
    XORL AX, AX
    MOVL AX, 276(SP)
    MOVL AX, 280(SP)
    MOVL AX, 284(SP)
    MOVL AX, 288(SP)
    MOVL AX, 292(SP)
    MOVL AX, 296(SP)
    MOVL AX, 300(SP)
    MOVL AX, 304(SP)
    MOVL AX, 308(SP)

    // sigfillset(&ign = &16(SP))
    LEAL 16(SP), AX
    MOVL AX, 0(SP)
    MOVL $sigfillset__dynload(SB), AX
    CALL (AX)

    // pthread_sigmask(SIG_SETMASK = 2, &ign = &16(SP), &oset = &148(SP))
    MOVL $2, 0(SP)
    LEAL 16(SP), AX
    MOVL AX, 4(SP)
    LEAL 148(SP), AX
    MOVL AX, 8(SP)
    MOVL $pthread_sigmask__dynload(SB), AX
    CALL (AX)

    // pthread_attr_init(&attr = &276(SP))
    LEAL 276(SP), AX
    MOVL AX, 0(SP)
    MOVL $pthread_attr_init__dynload(SB), AX
    CALL (AX)

    // pthread_attr_getstacksize(&attr = &276(SP), &size = &12(SP))
    LEAL 276(SP), AX
    MOVL AX, 0(SP)
    LEAL 12(SP), AX
    MOVL AX, 4(SP)
    MOVL $pthread_attr_getstacksize__dynload(SB), AX
    CALL (AX)

    // ts->g->stack_high = size
    MOVL ts+0(FP), AX
    MOVL 0(AX), AX
    MOVL 12(SP), BX
    MOVL BX, stack_hi(AX)

    // In the cgocode this has *pthread_t as first argument, but we don't use this in this function...
    // And by not doing that it is simpler here
    // SI = _cgo_try_pthread_create(&attr = 276(SP), threadentry, &ts)
    LEAL 276(SP), AX
    MOVL AX, 0(SP)
    MOVL $threadentry(SB), AX
    MOVL AX, 4(SP)
    MOVL ts+0(FP), AX
    MOVL AX, 8(SP)
    // We didn't clobber any registers here -> no saving required
    CALL _cgo_try_pthread_create(SB)

    MOVL AX, SI

    // pthread_sigmask(SIG_SETMASK = 2, &oset = &148(SP), nil)
    MOVL $2, 0(SP)
    LEAL 148(SP), AX
    MOVL AX, 4(SP)
    MOVL $0, 8(SP)
    MOVL $pthread_sigmask__dynload(SB), AX
    CALL (AX)

    // if R12 != 0 -> fail
    TESTL SI, SI
    JZ OK

    INT $3 // pthread_create failed

OK:
    RET

// _cgo_try_pthread_create(*attr, *threadentry, *ts) (runtime/cgo/gcc_libinit.c)
//  0(SP) arg0
//  4(SP) arg1
//  8(SP) arg2
// 12(SP) arg3
// 16(SP) pthread_t (4)
// 20(SP) timespec_t (8)
TEXT _cgo_try_pthread_create(SB),NOSPLIT,$28-12    
    MOVL $0, BX // tries

AGAIN:
    // AX = pthread_create(&pthread_t = &0(SP), attr, threadentry, ts)
    LEAL 16(SP), AX
    MOVL AX, 0(SP)
    MOVL attr+0(FP), AX
    MOVL AX, 4(SP)
    MOVL attr+4(FP), AX
    MOVL AX, 8(SP)
    MOVL attr+8(FP), AX
    MOVL AX, 12(SP)
    MOVL $pthread_create__dynload(SB), AX
    CALL (AX)

    TESTL AX, AX
    JNZ FAILED

    // if AX == 0
    // pthread_detach(pthread_t)
    MOVL 16(SP), AX
    MOVL AX, 0(SP)
    MOVL $pthread_detach__dynload(SB), AX
    CALL (AX)

    // return 0
    XORL AX, AX
    RET

    // AX != 0:
FAILED:
    CMPL AX, $11 //EAGAIN
    JEQ WAIT

    // AX != EAGAIN -> return the error
    RET

WAIT:
    // AX == EAGAIN

    // ts.tv_sec = 0
    MOVL $0, 20(SP)
    // ts.tv_nsec = (1+tries)*1000*1000
    MOVL BX, AX
    INCL AX
    IMULL $(1000*1000), AX
    MOVL AX, 24(SP)
    
    // nanosleep(&ts = &8(SP), nil)
    LEAL 20(SP), AX
    MOVL AX, 0(SP)
    XORL AX, AX
    MOVL AX, 4(SP)
    MOVL $nanosleep__dynload(SB), AX
    CALL (AX)

    // tries++
    INCL BX
    CMPL BX, $20
    JLT AGAIN

    // if tries >= 20 return eagain

    MOVL $11, AX //EAGAIN
    RET

// threadentry(ThreadStart *ts) (runtime/cgo/gcc_linux_amd64.c)
// 0(SP) arg0
// 4(SP) ThreadStart (size 12)
TEXT threadentry(SB),NOSPLIT,$16-4 //ThreadStart
    // 4(SP) = *ts
    MOVL ts+0(FP), BX
    MOVL 0(BX), AX
    MOVL AX, 4(SP)
    MOVL 4(BX), AX
    MOVL AX, 8(SP)
    MOVL 8(BX), AX
    MOVL AX, 12(SP)

    // free(ts)
    // ts is still the first argument
    MOVL BX, 0(SP)
    MOVL $free__dynload(SB), AX
    CALL (AX)

    // setg_gcc(ts->g)
    MOVL 4(SP), AX
    MOVL AX, 0(SP)
    MOVL $x_setg_gcc(SB), AX // maybe this is bad?
    CALL (AX)

    // ts.fn
    // in gcc_linux_amd64.c this uses crosscall, which we don't need here (we don't care about any register after this point)
    // also no need for clearing AX -> we call into go
    MOVL 12(SP), AX
    CALL AX

    // return 0
    XORL AX, AX
    RET

// _cgo_notify_runtime_init_done (runtime/cgo/gcc_libinit.c)
// do nothing - we don't support being a library
TEXT x_cgo_notify_runtime_init_done(SB),NOSPLIT,$0
    RET

// _cgo_setenv(char **arg) (runtime/cgo/gcc_setenv.c)
// 0(SP) arg0
// 4(SP) arg1
// 8(SP) arg2
TEXT x_cgo_setenv(SB),NOSPLIT,$12-4
    // setenv(arg[0], arg[1], 1)
    MOVL arg+0(FP), BX
    MOVL 0(BX), AX
    MOVL AX, 0(SP)
    MOVL 4(BX), AX
    MOVL AX, 4(SP)
    MOVL $1, 8(SP)
    MOVL $setenv__dynload(SB), AX
    CALL (AX)
    RET

// _cgo_unsetenv(char *arg) (runtime/cgo/gcc_setenv.c)
TEXT x_cgo_unsetenv(SB),NOSPLIT,$4-4
    // unsetenv(arg)
    MOVL arg+0(FP), AX
    MOVL AX, 0(SP)
    MOVL $setenv__dynload(SB), AX
    CALL (AX)
    RET

// _cgo_callers(uintptr_t sig, void *info, void *context, void (*cgoTraceback)(struct cgoTracebackArg*), uintptr_t* cgoCallers, void (*sigtramp)(uintptr_t, void*, void*)) (runtime/cgo/gcc_traceback.c)
//  0(SP) arg0
//  4(SP) arg1
//  8(SP) arg2
// 12(SP) cgoTracebackArg
TEXT x_cgo_callers(SB),NOSPLIT,$28
    // arg.Context = 0
    MOVL $0, 12(SP)
    // arg.SigContext = context
    MOVL context+8(FP), AX
    MOVL AX, 16(SP)
    // arg.Buf = cgoCallers
    MOVL cgoCallers+16(FP), AX
    MOVL AX, 20(SP)
    // arg.Max = 32
    MOVL $32, 24(SP)

    // cgoTraceback(&arg)
    MOVL cgoTraceback+12(FP), BX
    LEAL 12(SP), AX
    MOVL AX, 0(SP)
    CALL (BX)

    // sigtramp(sig, info, context)
    MOVL sig+0(FP), AX
    MOVL AX, 0(SP)
    MOVL info+4(FP), AX
    MOVL AX, 4(SP)
    MOVL context+8(FP), AX
    MOVL AX, 8(SP)
    MOVL sigtramp+20(FP), AX
    CALL (AX)
    RET
