// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin openbsd,amd64 openbsd,arm64

package forkexec

import (
    "syscall"
    "unsafe"
)

//go:linkname runtime_BeforeFork syscall.runtime_BeforeFork
func runtime_BeforeFork()

//go:linkname runtime_AfterFork syscall.runtime_AfterFork
func runtime_AfterFork()

//go:linkname runtime_AfterForkInChild syscall.runtime_AfterForkInChild
func runtime_AfterForkInChild()

//go:linkname rawSyscall syscall.rawSyscall
func rawSyscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:linkname libc_fork_trampoline syscall.libc_fork_trampoline
func libc_fork_trampoline()

//go:linkname libc_setsid_trampoline syscall.libc_setsid_trampoline
func libc_setsid_trampoline()

//go:linkname libc_setpgid_trampoline syscall.libc_setpgid_trampoline
func libc_setpgid_trampoline()

//go:linkname libc_getpid_trampoline syscall.libc_getpid_trampoline
func libc_getpid_trampoline()

//go:linkname libc_ioctl_trampoline syscall.libc_ioctl_trampoline
func libc_ioctl_trampoline()

//go:linkname libc_chroot_trampoline syscall.libc_chroot_trampoline
func libc_chroot_trampoline()

//go:linkname libc_setgroups_trampoline syscall.libc_setgroups_trampoline
func libc_setgroups_trampoline()

//go:linkname libc_setgid_trampoline syscall.libc_setgid_trampoline
func libc_setgid_trampoline()

//go:linkname libc_setuid_trampoline syscall.libc_setuid_trampoline
func libc_setuid_trampoline()

//go:linkname libc_chdir_trampoline syscall.libc_chdir_trampoline
func libc_chdir_trampoline()

//go:linkname libc_dup2_trampoline syscall.libc_dup2_trampoline
func libc_dup2_trampoline()

//go:linkname libc_fcntl_trampoline syscall.libc_fcntl_trampoline
func libc_fcntl_trampoline()

//go:linkname libc_close_trampoline syscall.libc_close_trampoline
func libc_close_trampoline()

//go:linkname libc_execve_trampoline syscall.libc_execve_trampoline
func libc_execve_trampoline()

//go:linkname libc_write_trampoline syscall.libc_write_trampoline
func libc_write_trampoline()

//go:linkname libc_exit_trampoline syscall.libc_exit_trampoline
func libc_exit_trampoline()

//go:linkname libc_read_trampoline syscall.libc_read_trampoline
func libc_read_trampoline()

//go:linkname fcntl syscall.fcntl
func fcntl(fd int, cmd int, arg int) (val int, err error)

//go:linkname ptrace1 syscall.ptrace1
func ptrace1(request int, pid int, addr uintptr, data uintptr) (err error)

//go:linkname readlen syscall.readlen
func readlen(fd int, buf *byte, nbuf int) (n int, err error)


// Find the entry point for f. See comments in runtime/proc.go for the
// function of the same name.
//go:nosplit
func funcPC(f func()) uintptr {
    return **(**uintptr)(unsafe.Pointer(&f))
}

// Fork, dup fd onto 0..len(fd), and exec(argv0, argvv, envv) in child.
// If a dup or exec fails, write the errno error to pipe.
// (Pipe is close-on-exec so if exec succeeds, it will be closed.)
// In the child, this function must not acquire any locks, because
// they might have been locked at the time of the fork. This means
// no rescheduling, no malloc calls, and no new stack segments.
// For the same reason compiler does not race instrument it.
// The calls to rawSyscall are okay because they are assembly
// functions that do not grow the stack.
//go:norace
func forkAndExecInChild(argv0 *byte, argv, envv []*byte, chroot, dir *byte, attr *syscall.ProcAttr, sys *syscall.SysProcAttr, pipe int) (pid int, err syscall.Errno) {
    // Declare all variables at top in case any
    // declarations require heap allocation (e.g., err1).
    var (
        r1     uintptr
        err1   syscall.Errno
        nextfd int
        i      int
    )

    // guard against side effects of shuffling fds below.
    // Make sure that nextfd is beyond any currently open files so
    // that we can't run the risk of overwriting any of them.
    fd := make([]int, len(attr.Files))
    nextfd = len(attr.Files)
    for i, ufd := range attr.Files {
        if nextfd < int(ufd) {
            nextfd = int(ufd)
        }
        fd[i] = int(ufd)
    }
    nextfd++

    // About to call fork.
    // No more allocation or calls of non-assembly functions.
    runtime_BeforeFork()
    r1, _, err1 = rawSyscall(funcPC(libc_fork_trampoline), 0, 0, 0)
    if err1 != 0 {
        runtime_AfterFork()
        return 0, err1
    }

    if r1 != 0 {
        // parent; return PID
        runtime_AfterFork()
        return int(r1), 0
    }

    // Fork succeeded, now in child.

    runtime_AfterForkInChild()

    // Enable tracing if requested.
    if sys.Ptrace {
        if err := ptrace(syscall.PTRACE_TRACEME, 0, 0, 0); err != nil {
            err1 = err.(syscall.Errno)
            goto childerror
        }
    }

    // Session ID
    if sys.Setsid {
        _, _, err1 = rawSyscall(funcPC(libc_setsid_trampoline), 0, 0, 0)
        if err1 != 0 {
            goto childerror
        }
    }

    // Set process group
    if sys.Setpgid || sys.Foreground {
        // Place child in process group.
        _, _, err1 = rawSyscall(funcPC(libc_setpgid_trampoline), 0, uintptr(sys.Pgid), 0)
        if err1 != 0 {
            goto childerror
        }
    }

    if sys.Foreground {
        pgrp := sys.Pgid
        if pgrp == 0 {
            r1, _, err1 = rawSyscall(funcPC(libc_getpid_trampoline), 0, 0, 0)
            if err1 != 0 {
                goto childerror
            }

            pgrp = int(r1)
        }

        // Place process group in foreground.
        _, _, err1 = rawSyscall(funcPC(libc_ioctl_trampoline), uintptr(sys.Ctty), uintptr(syscall.TIOCSPGRP), uintptr(unsafe.Pointer(&pgrp)))
        if err1 != 0 {
            goto childerror
        }
    }

    // Chroot
    if chroot != nil {
        _, _, err1 = rawSyscall(funcPC(libc_chroot_trampoline), uintptr(unsafe.Pointer(chroot)), 0, 0)
        if err1 != 0 {
            goto childerror
        }
    }

    // User and groups
    if cred := sys.Credential; cred != nil {
        ngroups := uintptr(len(cred.Groups))
        groups := uintptr(0)
        if ngroups > 0 {
            groups = uintptr(unsafe.Pointer(&cred.Groups[0]))
        }
        if !cred.NoSetGroups {
            _, _, err1 = rawSyscall(funcPC(libc_setgroups_trampoline), ngroups, groups, 0)
            if err1 != 0 {
                goto childerror
            }
        }
        _, _, err1 = rawSyscall(funcPC(libc_setgid_trampoline), uintptr(cred.Gid), 0, 0)
        if err1 != 0 {
            goto childerror
        }
        _, _, err1 = rawSyscall(funcPC(libc_setuid_trampoline), uintptr(cred.Uid), 0, 0)
        if err1 != 0 {
            goto childerror
        }
    }

    // Chdir
    if dir != nil {
        _, _, err1 = rawSyscall(funcPC(libc_chdir_trampoline), uintptr(unsafe.Pointer(dir)), 0, 0)
        if err1 != 0 {
            goto childerror
        }
    }

    // Pass 1: look for fd[i] < i and move those up above len(fd)
    // so that pass 2 won't stomp on an fd it needs later.
    if pipe < nextfd {
        _, _, err1 = rawSyscall(funcPC(libc_dup2_trampoline), uintptr(pipe), uintptr(nextfd), 0)
        if err1 != 0 {
            goto childerror
        }
        rawSyscall(funcPC(libc_fcntl_trampoline), uintptr(nextfd), syscall.F_SETFD, syscall.FD_CLOEXEC)
        pipe = nextfd
        nextfd++
    }
    for i = 0; i < len(fd); i++ {
        if fd[i] >= 0 && fd[i] < int(i) {
            if nextfd == pipe { // don't stomp on pipe
                nextfd++
            }
            _, _, err1 = rawSyscall(funcPC(libc_dup2_trampoline), uintptr(fd[i]), uintptr(nextfd), 0)
            if err1 != 0 {
                goto childerror
            }
            rawSyscall(funcPC(libc_fcntl_trampoline), uintptr(nextfd), syscall.F_SETFD, syscall.FD_CLOEXEC)
            fd[i] = nextfd
            nextfd++
        }
    }

    // Pass 2: dup fd[i] down onto i.
    for i = 0; i < len(fd); i++ {
        if fd[i] == -1 {
            rawSyscall(funcPC(libc_close_trampoline), uintptr(i), 0, 0)
            continue
        }
        if fd[i] == int(i) {
            // dup2(i, i) won't clear close-on-exec flag on Linux,
            // probably not elsewhere either.
            _, _, err1 = rawSyscall(funcPC(libc_fcntl_trampoline), uintptr(fd[i]), syscall.F_SETFD, 0)
            if err1 != 0 {
                goto childerror
            }
            continue
        }
        // The new fd is created NOT close-on-exec,
        // which is exactly what we want.
        _, _, err1 = rawSyscall(funcPC(libc_dup2_trampoline), uintptr(fd[i]), uintptr(i), 0)
        if err1 != 0 {
            goto childerror
        }
    }

    // By convention, we don't close-on-exec the fds we are
    // started with, so if len(fd) < 3, close 0, 1, 2 as needed.
    // Programs that know they inherit fds >= 3 will need
    // to set them close-on-exec.
    for i = len(fd); i < 3; i++ {
        rawSyscall(funcPC(libc_close_trampoline), uintptr(i), 0, 0)
    }

    // Detach fd 0 from tty
    if sys.Noctty {
        _, _, err1 = rawSyscall(funcPC(libc_ioctl_trampoline), 0, uintptr(syscall.TIOCNOTTY), 0)
        if err1 != 0 {
            goto childerror
        }
    }

    // Set the controlling TTY to Ctty
    if sys.Setctty {
        _, _, err1 = rawSyscall(funcPC(libc_ioctl_trampoline), uintptr(sys.Ctty), uintptr(syscall.TIOCSCTTY), 0)
        if err1 != 0 {
            goto childerror
        }
    }

    // Time to exec.
    _, _, err1 = rawSyscall(funcPC(libc_execve_trampoline),
        uintptr(unsafe.Pointer(argv0)),
        uintptr(unsafe.Pointer(&argv[0])),
        uintptr(unsafe.Pointer(&envv[0])))

childerror:
    // send error code on pipe
    rawSyscall(funcPC(libc_write_trampoline), uintptr(pipe), uintptr(unsafe.Pointer(&err1)), unsafe.Sizeof(err1))
    for {
        rawSyscall(funcPC(libc_exit_trampoline), 253, 0, 0)
    }
}

