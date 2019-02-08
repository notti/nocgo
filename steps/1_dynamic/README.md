ASM static binary -> dynamic loading without linking with `ld`
==============================================================

This directory contains the same 42.asm as 0_static, with an added call to `__printf` and argument "hello world" that will be replaced by `puts` from glibc.

`rewrite.go` takes `a.out` as input and outputs `dyn` with changed elf headers so that:
1. `interp` is added so we get loaded dynamically
2. `libc.so.6` is added as needed library
3. contents of `__printf` in `.bss` get changed to point to `puts` from `libc.so.6`
4. profit: `./dyn` outputs `hello world`, is statically compiled (without libc!), but can call into the dynamic symbol `puts`