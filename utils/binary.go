package utils

import (
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"syscall"
)

func IsExecutableFile (filePath string) (bool, error) {
	fp, err := os.OpenFile(filePath, os.O_RDONLY | syscall.O_NONBLOCK, 0)
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


