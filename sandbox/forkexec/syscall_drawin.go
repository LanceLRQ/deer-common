// Copyright 2009,2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Darwin system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and wrap
// it in our own nicer implementation, either here or in
// syscall_bsd.go or syscall_unix.go.

package forkexec

import "unsafe"

func readlen(fd int, buf *byte, nbuf int) (n int, err error) {
    r0, _, e1 := rsyscall(funcPC(libc_read_trampoline), uintptr(fd), uintptr(unsafe.Pointer(buf)), uintptr(nbuf))
    n = int(r0)
    if e1 != 0 {
        err = errnoErr(e1)
    }
    return
}
