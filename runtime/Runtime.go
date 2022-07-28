package runtime

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// when AlwaysStackTrace is true, Caller() will always return CallerStackTrace() result
var AlwaysStackTrace = false

// 取得呼叫的文件與行號
func Caller(skip int) string {
	if AlwaysStackTrace {
		return CallerStackTrace(skip + 1)
	}
	_, file, line, ok := runtime.Caller(1 + skip)
	if !ok {
		return "[fail to get caller]"
	}

	return fmt.Sprintf("%s/%s:%d", filepath.Base(filepath.Dir(file)), filepath.Base(file), line)
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
