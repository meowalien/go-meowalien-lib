package sign

import (
	"crypto/md5"
	"encoding/hex"
	"log"
)

func Md5Sign(bt []byte) string {
	md5Ctx := md5.New()
	_, err := md5Ctx.Write(bt)
	if err != nil {
		log.Println("error when Write to md5 Hash: ", err.Error())
		return ""
	}
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
