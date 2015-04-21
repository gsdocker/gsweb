package gsweb

import "os"

func fileIsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func isDir(filename string) bool {

	fi, err := os.Stat(filename)

	if err != nil {
		return false
	}
	return fi.IsDir()
}
