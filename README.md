nocgo
=====

At the moment only go1.11 is supported. go1.12 causes some problems with nosplit.

[![GoDoc](https://godoc.org/github.com/notti/nocgo?status.svg)](https://godoc.org/github.com/notti/nocgo)

This repository/package contains a *proof of concept* for calling into C code *without* using cgo.

> **WARNING!** This is meant as a proof of concept and subject to changes.
Furthermore this is highly experimental code. DO NOT USE IN PRODUCTION.
This could cause lots of issues from random crashes (there are tests - but there is definitely stuff that's not tested) to teaching your gopher to talk [C gibberish](https://cdecl.org/).

> **WARNING** nocgo supports both cgo and missing cgo as environment. So if you want to ensure cgo not being used don't forget `CGO_ENABLED=0` as environment variable to `go build`.

Usage
-----

Libraries can be loaded and unloaded similar to `dlopen` and `dlclose`, but acquiring symbols (i.e., functions, global variables) is a bit different, since a function specification (i.e., arguments, types, return type) is also needed. Furthermore, C-types must be translated to go-types and vice versa.

This works by providing a function specification as a `struct`, where all the elements are the function arguments (in the same order) and one element can be marked as return (with the tag `nocgo:"ret"`). Before function call, the argument values must be provided with this `struct` and after the function call, the return value will be in the return element. The function call might also change argument values, if those are provided via pointers.

### Type Mappings

Go types will be mapped to C-types according to the following table:

Go type                                       | C Type
--------------------------------------------- | ------
`int8`, `byte`                                | `char`
`uint8`, `bool`                               | `unsigned char`
`int16`                                       | `short`
`uint16`                                      | `unsigned short`
`int32`                                       | `int`
`uint32`                                      | `unsigned int`
`int64`                                       | `long`
`uint64`                                      | `unsigned long`
`float32`                                     | `float`
`float64`                                     | `double`
`[]`, `uintptr`, `reflect.UnsafePointer`, `*` | `*`

The last line means that slices and pointers are mapped to pointers in C. Pointers to structs are possible.

Passing `struct`, `complex`, and callback functions is not (yet) supported.

> **WARNING** `struct`s that are referenced **must** follow C alignment rules! There is **no** type checking, since this is actually not possible due to libraries not knowing their types...

Go `int` was deliberately left out to avoid confusion, since it has different sizes on different architectures.

### Example

An example using `pcap_open_live` from libpcap (C-definition: `pcap_t *pcap_open_live(const char *device, int snaplen, int promisc, int to_ms, char *errbuf)
`) could look like the following example:

```golang
// Argument specification
type pcapOpenLiveArgs struct {
    device  []byte
    snaplen int32
    promisc int32
    toMS    int32
    errbuf  []byte
    ret     uintptr `nocgo:"ret"`
}
// Provide some values
openLiveArg := pcapOpenLiveArgs{
    device:  nocgo.MakeCString("lo"),
    snaplen: 1500,
    promisc: 1,
    toMS:    100,
    errbuf:  make([]byte, 512),
}
// Load the library
lib, err := nocgo.Open("libpcap.so")
if err != nil {
    log.Fatalln("Couldn't load libpcap: ", err)
}
// Get a handle for the function
pcapOpenLive, err := lib.Func("pcap_open_live", openLiveArg) // here we could also provide pcapOpenLiveArgs{}
if err != nil {
    log.Fatalln("Couldn't get pcap_open_live: ", err)
}

// Do the function call
pcapOpenLive.Call(unsafe.Pointer(&openLiveArg))

// Check return value
if openLiveArg.ret == 0 {
    log.Fatalf("Couldn't open %s: %s\n", *dev, nocgo.MakeGoStringFromSlice(openLiveArg.errbuf))
}

// openLiveArg.ret can now be used as argument to the other libpcap functions
```

A full example is contained in [examplelibpcap](examplelibpcap) and another one in [example](example).

> **WARNING** nocgo supports both cgo and missing cgo as environment. So if you want to ensure cgo not being used don't forget `CGO_ENABLED=0` as environment variable to `go build`.

Supported Systems
-----------------

* linux with glibc
* FreeBSD<br>
  *Errata:* FreeBSD requires the exported symbols `_environ` and `_progname`. This is only possible inside cgo or stdlib. So for building on FreeBSD, `-gcflags=github.com/notti/nocgo/fakecgo=-std` is required (This doesn't seem to work for `go test` - so examples work, but test does not)).

With some small modifications probably all systems providing `dlopen` can be supported. Have a look at [dlopen_OS.go](dlopen_linux.go) and [symbols_OS.go](fakecgo/symbols_linux.go) in fakecgo.

Supported Architectures
-----------------------

* 386
* amd64

Implementing further architectures requires
* Building trampolines for [fakecgo](fakecgo) (see below)
* Implementing the cdecl callspec in [call_.go](call_amd64.go)/[.s](call_amd64.s)

How does this work
------------------

### nocgo

nocgo imports `dlopen`, `dlclose`, `dlerror`, `dlsym` via `go:cgo_import_dynamic` in [dlopen_OS.go](dlopen_linux.go). `lib.Func` builds a specification on where to put which argument in [call_arch.go](call_amd64.go). `spec.Call` uses `cgocall` from the runtime to call an assembly function and pass the spec to it. This assembly function is implemented in call_arch.s and it uses the specification to place the arguments into the right places, calls the pointer provided by `dlsym` and then puts the return argument into the right place if needed.

This is basically what `libffi` does. So far cdecl for 386 (pass arguments on the stack in right to left order, return values are in AX/CX or ST0) and amd64 (pass arguments in registers DI, SI, DX, CX, R8, R9/X0-X7 and the stack in right to left order, number of floats in AX, fixup alignment of stack) are implemented.

So far so simple. `cgocall` could actually be used to call a C function directly - but it is only capable of providing one argument!

But there is a second issue. For simple C functions we could leave it at that (well we would need to use `asmcgocall`, because `cgocall` checks, if cgo is actually there...). But there is this thing called Thread Local Storage (TLS) that is not too happy about golang not setting that up correctly. This is already needed if you do `printf("%f", 1)` with glibc!

So we need to provide some functionality that cgo normally provides, which is implemented in fakecgo:

### fakecgo

go sets up it's own TLS during startup in runtime/asm_arch.s in `runtime·rt0_go`. We can easily prevent that by providing setting the global variable `_cgo_init` to something non-zero (easily achieved with `go:linkname` and setting a value). But this would crash go, since if this is the case, go actually calls the address inside this variable (well ok we can provide an empty function).

Additionally, this would provide correct TLS only on the main thread. This works until one does a lot more than just call one function, so we need to fixup also some other stuff.

So next step: set `runtime.is_cgo` to true (again - linkname to the rescue). But this will panic since now the runtime expects the global variables `_cgo_thread_start`, `_cgo_notify_runtime_init_done`, `_cgo_setenv`, and `_cgo_unsetenv` to point to something. Ok so let's just implement those.

* `_cgo_notify_runtime_init_done` is easy - we don't need this one: empty function.
* `_cgo_setenv` is also simple: just one function call to `setenv`
* `_cgo_unsetenv` is the same.
* `_cgo_init` queries the needed stack size to update g->stack so that runtime stack checks do the right thing (it also provides a setg function we come to that later...)
* `_cgo_thread_start` is a bit more involved... It starts up a new thread with `pthread_create` and does a bit of setup.

So this should be doable - right?

Well easier said than done - those are implemented in C-code in runtime/cgo/*c presenting some kind of chicken and egg problem to us.

So I started out with reimplementing those in go assembly (remember: we want to be cgo free) which is available in the tag asm. Since this is really cumbersome and needs a lot of code duplication, I experimented a bit if we can do better.

Aaaand we can:

[fakecgo/trampoline_arch.s](fakecgo/trampoline_amd64.s) contains the above mentioned entry points, and "converts" the C-calling conventions to go calling conventions (e.g. move register passed arguments to the stack). Then it calls the go functions in [fakecgo/cgo.go](fakecgo/cgo.go).

Ok - but we still need all those pthread and C-library-functions. Well we can import the symbols (like with `dlopen`). So all we need is a way to call those:

The trampoline file also contains an `asmlibccall6` function that can call C-functions with a maximum of 6 integer arguments and one return value. [fakecgo/libccall.go](fakecgo/libccall.go) maps this onto more convenient go functions with 1-6 arguments and [fakecgo/libcdefs.go](fakecgo/libcdefs.go) further maps those into nice functions that look like the C functions (e.g. `func pthread_create(thread *pthread_t, attr *pthread_attr, start, arg unsafe.Pointer) int32`). Well this was not exactly my idea - the runtime already does that for solaris and darwin (runtime/os_solaris.go, runtime/syscall_solaris.go, runtime/sys_solaris_amd64.s) - although my implementation here is kept a bit simpler since it only ever will be called from gocode pretending to be C.

So now we can implement all the above mentioned cgo functions in pure (but sometimes a bit ugly) go in [fakecgo/cgo.go](fakecgo/cgo.go). Ugly, because those functions are called with lots of functionality missing! Writebarriers are **not** allowed, as are stack splits.

The upside is, that the only arch dependent stuff are the trampolines (in assembly) and the only OS dependent stuff are the symbol imports.

Except for freebsd (which needs two exported symbols, as mentioned above) all those things work outside the runtime and no special treatment is needed. Just import fakecgo and all the cgo setup just works (except if you use cgo at the same time - then the linker will complain).

Benchmarks
----------

This will be a bit slower than cgo. Most of this is caused by argument rearranging:

### 386

```
name           old time/op    new time/op    delta
Empty-4          89.8ns ±12%    85.6ns ± 5%      ~     (p=0.481 n=10+10)
Float2-4         84.6ns ± 1%   215.8ns ± 2%  +154.91%  (p=0.000 n=8+9)
StackSpill3-4     118ns ± 5%     126ns ± 5%    +7.07%  (p=0.000 n=10+8)
```

Float is so slow since that type is at the end of the comparison chain.

### amd64

```
name           old time/op    new time/op    delta
Empty-4          70.1ns ± 8%    73.9ns ± 3%   +5.36%  (p=0.026 n=10+9)
Float2-4         72.0ns ± 4%    90.9ns ± 4%  +26.20%  (p=0.000 n=10+10)
StackSpill3-4    88.5ns ± 4%   117.2ns ± 1%  +32.52%  (p=0.000 n=10+8)
```
