package redisdb

import (
	"fmt"
	"strings"
)

var Split = ":"
var Prefix = ""

// 格式化redis key
func FormatKey(keys ...interface{}) string {
	l := len(keys)

	var k = make([]string, l+1)

	k[0] = Prefix
	for i, key := range keys {
		switch kt := key.(type) {
		case string:
			k[i+1] = kt
		default:
			k[i+1] = fmt.Sprint(kt)
		}
	}

	return strings.Join(k, Split)
}
