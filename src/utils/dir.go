package utils

import (
	"errors"
	"os"
)

var (
	ErrArgNotDir = errors.New("Argument is not a directory.")
)

func FilesUnderDir(dir string, num int) (files []string, err error) {
	dirFile, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer dirFile.Close()
	dirFileStat, err := dirFile.Stat()
	if err != nil {
		return nil, err
	}
	if !dirFileStat.IsDir() {
		return nil, ErrArgNotDir
	}
	files, err = dirFile.Readdirnames(num)
	return files, err
}

func AllFilesUnderDir(dir string) (files []string, err error) {
	files, err = FilesUnderDir(dir, -1)
	return
}
