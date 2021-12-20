package hash

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
)

func NewSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func NewMD5(data []byte) []byte {
	h := md5.New()
	_ , err  := bufio.NewWriter(h).Write(data)
	if err != nil{
		return nil
	}
	return h.Sum(nil)
}
