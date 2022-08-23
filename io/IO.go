package io

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"os"
)

// 將檔案存到硬碟
func WriteToNewFile(fp string, bin []byte) (err error) {
	fo, err := os.Create(fp)
	if err != nil {
		return errs.New("error when os.create: %w", err)
	}
	defer func(fo *os.File) {
		e := fo.Close()
		if e != nil {
			fmt.Printf("error when close file: %s", e.Error())
		}
	}(fo)

	_, err = fo.Write(bin)
	if err != nil {
		return errs.New("error Write file: %w", err)
	}
	return nil
}
