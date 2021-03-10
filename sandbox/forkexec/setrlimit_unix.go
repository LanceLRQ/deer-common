// +build darwin linux

package forkexec

import (
    "syscall"
)

// 定义ITimer的常量
const (
    ITIMER_REAL    = 0
    ITIMER_VIRTUAL = 1
    ITIMER_PROF    = 2
)

type RLimit struct {
    Which  int
    Enable bool
    RLim   syscall.Rlimit
}

type ITimerVal struct {
    ItInterval TimeVal
    ItValue    TimeVal
}

type TimeVal struct {
    TvSec  uint64
    TvUsec uint64
}

type ForkExecRLimit struct {
    TimeLimit     int
    RealTimeLimit int
    MemoryLimit   int
    FileSizeLimit int
    StackLimit    int
}
