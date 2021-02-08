package problems

import (
    "archive/zip"
    "bufio"
    "fmt"
    "golang.org/x/crypto/openpgp"
)

// 读取题目携带的GPG信息
func ReadProblemGPGInfoZip(problemFile string) (string, error) {
    zipReader, err := zip.OpenReader(problemFile)
    if err != nil {
        return "", err
    }
    defer zipReader.Close()

    file, err := FindInZip(zipReader, ".gpg")
    if err != nil {
        if IsFileNotFoundError(err) {
            return "", fmt.Errorf("no GPG public key")
        }
        return "", err
    }

    elist, err := openpgp.ReadArmoredKeyRing(bufio.NewReader(*file))
    if err != nil {
        return "", err
    }
    if len(elist) < 1 {
        return "", fmt.Errorf("GPG key error")
    }
    rel := ""
    for _, identify := range elist[0].Identities {
        rel += identify.Name + "\n"
    }
    return rel, nil
}