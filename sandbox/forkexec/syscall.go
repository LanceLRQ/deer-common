package forkexec

import "syscall"

// Single-word zero for use when we need a valid pointer to 0 bytes.
// See mksyscall.pl.
var _zero uintptr


// 获取管道数据
func GetPipe() ([]uintptr, error) {
    var pipe = []int{0, 0}
    err := syscall.Pipe(pipe)
    if err != nil {
        return nil, err
    }
    return []uintptr{ uintptr(pipe[0]), uintptr(pipe[1]) }, nil
}