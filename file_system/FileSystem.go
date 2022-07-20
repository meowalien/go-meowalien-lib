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
	_, err := os.Stat(filename)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}
