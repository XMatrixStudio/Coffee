package services

import (
	"io"
	"os"
	"strings"
)

func copyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

// 判断文件夹是否存在
func pathExistsAndCreate(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	return nil
}

func indexOfFile(s []os.FileInfo, str string) int {
	for i := range s {
		name := strings.Split(s[i].Name(), "-")
		if len(name) < 2 {
			return -1
		}
		if name[1] == str+".png" {
			return i
		}
	}
	return -1
}
