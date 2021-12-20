package file_system

import (
	"errors"
	"io/fs"
	"os"
)

func FileMod(filename string) (fs.FileMode, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fi.Mode(), nil
}

func FileExist(filename string) (bool, error) {
	if _, err := os.Stat(filename); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
