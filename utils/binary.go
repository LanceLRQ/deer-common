package utils

import (
    "encoding/binary"
    "fmt"
    "github.com/LanceLRQ/deer-common/constants"
    "os"
    "path"
    "path/filepath"
    "runtime"
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

func GetCompiledBinaryFileAbsPath(typeName, moduleName, configDir string) (string, error) {
    targetName := GetCompiledBinaryFileName(typeName, moduleName)
    return filepath.Abs(path.Join(configDir, targetName))
}