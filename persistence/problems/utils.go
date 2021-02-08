package problems

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

func UnZip(zipArchive *zip.ReadCloser, destDir string) error {
    return WalkZip(zipArchive, func(f *zip.File) error {
        fpath := filepath.Join(destDir, f.Name)
        if f.FileInfo().IsDir() {
            _ = os.MkdirAll(fpath, os.ModePerm)
        } else {
            if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
                return err
            }

            inFile, err := f.Open()
            if err != nil {
                return err
            }
            defer inFile.Close()

            outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer outFile.Close()

            _, err = io.Copy(outFile, inFile)
            if err != nil {
                return err
            }
        }
        return nil
    })
}

func WalkZip(zipArchive *zip.ReadCloser, walkFunc func(file *zip.File) error) error {
    for _, f := range zipArchive.File {
        err := walkFunc(f)
        if err != nil {
            return err
        }
    }
    return nil
}

// 搜索zip内的文件并打开(精确匹配)
func FindInZip(zipArchive *zip.ReadCloser, fileName string) (*io.ReadCloser, error) {
    var fileResult io.ReadCloser
    finded := false
    err := walkZip(zipArchive, func(file *zip.File) error {
        if fileName == file.Name {
            finded = true
            var err error
            fileResult, err = file.Open()
            if err != nil {
                return err
            }
        }
        return nil
    })
    if !finded {
        return nil, fmt.Errorf("file (%s) not found", fileName)
    }
    return &fileResult, err
}