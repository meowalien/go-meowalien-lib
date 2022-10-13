package random

import (
	"encoding/binary"
	"github.com/bwmarrin/snowflake"
	"lukechampine.com/frand"
	"unsafe"
)

var snowflakeNode *snowflake.Node

// 節點編號
var NodeNum int64 = 1

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randInt64() int64 {
	b := make([]byte, 8)
	_, _ = frand.Read(b)
	return int64(binary.LittleEndian.Uint64(b))
}

func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randInt64(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randInt64(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b)) //nolint:gosec
}
