// +build linux,amd64

package forkexec

import "syscall"



//sysnb pipe2(p *[2]_C_int, flags int) (err error)

func Pipe2(p []int, flags int) (err error) {
    if len(p) != 2 {
        return syscall.EINVAL
    }
    var pp [2]_C_int
    err = pipe2(&pp, flags)
    p[0] = int(pp[0])
    p[1] = int(pp[1])
    return
}

