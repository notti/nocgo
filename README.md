nocgo
=====

Tested on go1.11 and go1.12.

[![GoDoc](https://godoc.org/github.com/notti/nocgo?status.svg)](https://godoc.org/github.com/notti/nocgo)

This repository/package contains a *proof of concept* for calling into C code *without* using cgo.

> **WARNING!** This is meant as a proof of concept and subject to changes.
Furthermore this is highly experimental code. DO NOT USE IN PRODUCTION.
This could cause lots of issues from random crashes (there are tests - but there is definitely stuff that's not tested) to teaching your gopher to talk [C gibberish](https://cdecl.org/).

> **WARNING** nocgo supports both cgo and missing cgo as environment. So if you want to ensure cgo not being used don't forget `CGO_ENABLED=0` as environment variable to `go build`.

Todo
----

- Callbacks into go
- Structures

When that's done write up a proposal for golang inclusion.

Usage
-----

Libraries can be loaded and unloaded similar to `dlopen` and `dlclose`, but acquiring symbols (i.e., functions, global variables) is a bit different, since a function specification (i.e., arguments, types, return type) is also needed. Furthermore, C-types must be translated to go-types and vice versa.

This works by providing a function specification as a pointer to a function variable. A call to `lib.Func` will examine arguments and eventual return value (only one or no return values allowed!), and set the function variable to a wrapper that will call into the desired C-function.

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

// Load the library
lib, err := nocgo.Open("libpcap.so")
if err != nil {
    log.Fatalln("Couldn't load libpcap: ", err)
}

// func specification
var pcapOpenLive func(device []byte, snaplen int32, promisc int32, toMS int32, errbuf []byte) uintptr
// Get a handle for the function
if err := lib.Func("pcap_open_live", &pcapOpenLive); err != nil {
    log.Fatalln("Couldn't get pcap_open_live: ", err)
}

// Do the function call
errbuf := make([]byte, 512)
pcapHandle := pcapOpenLive(nocgo.MakeCString("lo"), 1500, 1, 100, errbuf)

// Check return value
if pcapHandle == 0 {
    log.Fatalf("Couldn't open %s: %s\n", "lo", nocgo.MakeGoStringFromSlice(errbuf))
}

// pcapHandle can now be used as argument to the other libpcap functions
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

nocgo imports `dlopen`, `dlclose`, `dlerror`, `dlsym` via `go:cgo_import_dynamic` in [dlopen_OS.go](dlopen_linux.go). `lib.Func` builds a specification on where to put which argument in [call_arch.go](call_amd64.go). go calls such a function by dereferencing, where it points to, provide this address in a register and call the first address that is stored there. nocgo uses this mechanism by putting a struct there, that contains the address to a wrapper followed by a pointer to the what `dlsym` provided and a calling specification. The provided wrapper uses `cgocall` from the runtime to call an assembly function and pass the spec and a pointer to the arguments to it. This assembly function is implemented in call_arch.s and it uses the specification to place the arguments into the right places, calls the pointer provided by `dlsym` and then puts the return argument into the right place if needed.

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
Empty-4          84.5ns ± 0%    86.4ns ± 2%    +2.22%  (p=0.000 n=8+8)
Float2-4         87.9ns ± 1%   222.5ns ± 6%  +153.20%  (p=0.000 n=8+10)
StackSpill3-4     116ns ± 1%     130ns ± 1%   +12.04%  (p=0.000 n=8+8)
```

Float is so slow since that type is at the end of the comparison chain.

### amd64

```
name           old time/op    new time/op    delta
Empty-4          76.8ns ±10%    80.1ns ± 9%   +4.24%  (p=0.041 n=10+10)
Float2-4         78.4ns ± 5%    81.4ns ± 9%   +3.80%  (p=0.033 n=9+10)
StackSpill3-4    96.2ns ± 5%   120.7ns ± 7%  +25.46%  (p=0.000 n=10+9)
```
