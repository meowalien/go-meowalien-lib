package runtime

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type formatSetting int

const (
	CALLER_FORMAT_SHORT formatSetting = iota
	CALLER_FORMAT_LONG
)

// 取得呼叫的文件與行號
func Caller(skip int, formatSet formatSetting) string {
	_, file, line, ok := runtime.Caller(1 + skip)
	if !ok {
		return "[fail to get caller]"
	}
	switch formatSet {
	case CALLER_FORMAT_SHORT:
		dir, f := filepath.Split(file)
		return fmt.Sprintf("%s/%s:%d", filepath.Base(dir), f, line)
	case CALLER_FORMAT_LONG:
		return fmt.Sprintf("%s:%d", file, line)
	default:
		panic("unknown formatSetting")
	}
}

func CallerStackTrace(skip int) (ans string) {
	buf := make([]byte, 1024)
	var n int
	for {
		n = runtime.Stack(buf, false)
		if n < len(buf) {
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return cutOffStack(string(buf[:n]), skip)
}

func cutOffStack(ans string, skip int) string {
	i := strings.Index(ans, "\n")
	for skip += 1; skip > 0; skip-- {
		i += strings.Index(ans[i+1:], "\n") + 1
		i += strings.Index(ans[i+1:], "\n") + 1
	}

	return ans[i+1:]
}
