package redis

import (
	"fmt"
	"go.uber.org/zap/buffer"
)

const DefaultSplit = ":"
const DefaultPrefix = "root"

type KeyFormatter struct {
	Split  string
	Prefix string
}

// 格式化redis key
func (k KeyFormatter) Format(keys ...interface{}) string {
	bf := buffer.Buffer{}
	bf.AppendString(k.Prefix)
	for _, key := range keys {
		bf.AppendString(k.Split)
		switch kt := key.(type) {
		case string:
			bf.AppendString(kt)
		default:
			bf.AppendString(fmt.Sprint(kt))
		}
	}

	return bf.String()
}
