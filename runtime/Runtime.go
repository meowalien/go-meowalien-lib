package runtime

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// 取得呼叫的文件與行號
func CallerFileAndLine(deap int) string {
	_, file, line, ok := runtime.Caller(1 + deap)
	if !ok {
		return "[fail to get caller]"
	}

	return fmt.Sprintf("%s/%s:%d", filepath.Base(filepath.Dir(file)), filepath.Base(file), line)
}
