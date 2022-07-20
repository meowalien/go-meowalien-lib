package random

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
	"unsafe"
)

var snowflakeNode *snowflake.Node

// 節點編號
var NodeNum int64 = 1

func init() {
	rand.Seed(time.Now().UnixNano())
	var err error
	snowflakeNode, err = snowflake.NewNode(NodeNum)
	if err != nil {
		panic(err.Error())
	}
}

func Snowflake() snowflake.ID {
	return snowflakeNode.Generate()
}

func RandInt64(min int64, max int64) int64 {
	return min + rand.Int63n(max-min)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
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
