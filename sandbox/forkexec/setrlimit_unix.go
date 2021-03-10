// +build darwin linux

package forkexec

import (
    "math"
    "runtime"
    "syscall"
)

const DarwinSafeStackSize = 65500

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

type ExecRLimit struct {
    TimeLimit     int
    RealTimeLimit int
    MemoryLimit   int
    FileSizeLimit int
    StackLimit    int
}

type RlimitOptions struct {
    Rlimits     []RLimit
    ITimerValue ITimerVal
}

func GetRlimitOptions (sysRlimit *ExecRLimit) *RlimitOptions {
    // Make stack limit
    stackLimit := uint64(sysRlimit.StackLimit)
    if stackLimit <= 0 {
        stackLimit = uint64(sysRlimit.MemoryLimit * 1024)
    }
    if runtime.GOOS == "darwin" {
        if sysRlimit.MemoryLimit > DarwinSafeStackSize { // WTF?! >= 65mb caused an operation not permitted!
            stackLimit = uint64(DarwinSafeStackSize * 1024)
        }
    }

    return &RlimitOptions {
        Rlimits: []RLimit{
            // Set time limit: RLIMIT_CPU
            {
                Which: syscall.RLIMIT_CPU,
                Enable: sysRlimit.TimeLimit > 0,
                RLim: syscall.Rlimit{
                    Cur: uint64(math.Ceil(float64(sysRlimit.TimeLimit) / 1000.0)),
                    Max: uint64(math.Ceil(float64(sysRlimit.TimeLimit) / 1000.0)),
                },
            },
            // Set memory limit: RLIMIT_DATA
            {
                Which: syscall.RLIMIT_DATA,
                Enable: sysRlimit.MemoryLimit > 0,
                RLim: syscall.Rlimit{
                    Cur: uint64(sysRlimit.MemoryLimit * 1024),
                    Max: uint64(sysRlimit.MemoryLimit * 1024),
                },
            },
            // Set memory limit: RLIMIT_AS
            {
                Which: syscall.RLIMIT_AS,
                Enable: sysRlimit.MemoryLimit > 0,
                RLim: syscall.Rlimit{
                    Cur: uint64(sysRlimit.MemoryLimit * 1024 * 2),
                    Max: uint64(sysRlimit.MemoryLimit*1024*2 + 1024),
                },
            },
            // Set stack limit
            {
                Which: syscall.RLIMIT_STACK,
                Enable: stackLimit > 0,
                RLim: syscall.Rlimit{
                    Cur: stackLimit,
                    Max: stackLimit,
                },
            },
            // Set file size limit: RLIMIT_FSIZE
            {
                Which: syscall.RLIMIT_FSIZE,
                Enable: sysRlimit.FileSizeLimit > 0,
                RLim: syscall.Rlimit{
                    Cur: uint64(sysRlimit.FileSizeLimit),
                    Max: uint64(sysRlimit.FileSizeLimit),
                },
            },
        },
        ITimerValue: ITimerVal{
            ItInterval: TimeVal {
                TvSec: uint64(math.Floor(float64(sysRlimit.RealTimeLimit) / 1000.0)),
                TvUsec: uint64(sysRlimit.RealTimeLimit % 1000 * 1000),
            },
            ItValue: TimeVal{
                TvSec: uint64(math.Floor(float64(sysRlimit.RealTimeLimit) / 1000.0)),
                TvUsec: uint64(sysRlimit.RealTimeLimit % 1000 * 1000),
            },
        },
    }
}
