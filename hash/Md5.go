package hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap/buffer"
	"sort"
)

func Md5Hash(b []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(b)
	return md5Ctx.Sum(nil)
}
func Md5HashString(b []byte) string {
	return hex.EncodeToString(Md5Hash(b))
}

var bfPool = buffer.NewPool()

func Md5HashStringInUrlEncode(m map[string]interface{}, secretKey string, secret string, skipKey ...string) string {
	var str []string
loop:
	for k := range m {
		for _, s := range skipKey {
			if s == k {
				continue loop
			}
		}
		str = append(str, k)
	}
	if str == nil{return ""}

	sort.Strings(str)

	signstr := bfPool.Get()
	defer signstr.Free()
	defer signstr.Reset()
	for _, k := range str {
		if m == nil {
			continue
		}
		vv, ok := m[k]
		if !ok {
			continue
		}
		v := fmt.Sprintf("%v", vv)
		if v != "" {
			signstr.AppendString(k)
			signstr.AppendString("=")
			signstr.AppendString(v)
			signstr.AppendString("&")
		}
	}
	signstr.AppendString(secretKey)
	signstr.AppendString("=")
	signstr.AppendString(secret)
	return Md5HashString(signstr.Bytes())
}
