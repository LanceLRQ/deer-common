package utils

import (
    "bytes"
    "context"
    "encoding/binary"
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "github.com/LanceLRQ/deer-common/structs"
    "io"
    "os"
    "os/exec"
    "path"
    "path/filepath"
    "runtime"
    "strings"
    "syscall"
)

func IsExecutableFile(filePath string) (bool, error) {
    fp, err := os.OpenFile(filePath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
    if err != nil {
        return false, fmt.Errorf("open file error")
    }
    defer fp.Close()

    var magic uint32 = 0
    err = binary.Read(fp, binary.BigEndian, &magic)
    if err != nil {
        return false, err
    }

    isExec := false
    if runtime.GOOS == "darwin" {
        isExec = magic == 0xCFFAEDFE || magic == 0xCEFAEDFE || magic == 0xFEEDFACF || magic == 0xFEEDFACE
    } else if runtime.GOOS == "linux" {
        isExec = magic == 0x7F454C46
    }
    return isExec, nil
}

func GetCompiledBinaryFileName(typeName, moduleName string) string {
    prefix, ok := constants.TestlibBinaryPrefixs[typeName]
    if !ok { prefix = "" }
    return prefix + moduleName
}

// 根据配置文件将对应预编译文件转换成绝对路径
func GetCompiledBinaryFileAbsPath(typeName, moduleName, configDir string) (string, error) {
    targetName := GetCompiledBinaryFileName(typeName, moduleName)
    return filepath.Abs(path.Join(path.Join(configDir, "bin"), targetName))
}

// 解析generator脚本
func ParseGeneratorScript(script string) (string, []string, error) {
    vals := strings.Split(script, " ")
    if len(vals) <= 1 {
        return "", nil, fmt.Errorf("generator calling script error")
    }
    return vals[0], vals[1:], nil
}

// 运行UnixShell，支持context
func RunUnixShell(context context.Context, name string, args []string, onStart func(io.Writer) error) (*structs.ShellResult, error) {
    fpath, err := exec.LookPath(name)
    if err != nil {
        return nil, err
    }
    result := structs.ShellResult{}
    proc := exec.CommandContext(context, fpath, args...)
    var stderr, stdout bytes.Buffer
    proc.Stderr = &stderr
    proc.Stdout = &stdout
    stdin, err := proc.StdinPipe()
    if err != nil {
        return nil, err
    }
    //err = proc.Run()
    if err := proc.Start(); err != nil {
        return nil, err
    }
    if onStart != nil {
        err = onStart(stdin)
        if err != nil {
            return nil, err
        }
    }
    _ = stdin.Close()
    err = proc.Wait()
    result.Stdout = stdout.String()
    result.Stderr = stderr.String()
    result.ExitCode = proc.ProcessState.ExitCode()
    if err != nil {
        result.Success = false
        if serr := result.Stderr; serr == "" {
        //    result.Stderr += "\n" + err.Error()
        //} else {
            result.Stderr += err.Error()
        }
        return &result, nil
    }
    result.Success = true
    return &result, nil
}