package errs

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/runtime"
	"strings"
)

func WithLine(err interface{}, obj ...interface{}) error {
	switch errTp := err.(type) {
	case error:
		if len(obj) == 0 {
			return newWithLineError(runtime.CallerFileAndLine(1), errTp)
		} else if len(obj) == 1 && obj[0] != nil {
			if obj0, ok := obj[0].(error); ok {
				return wrapErr(errTp, obj0)
			}
		}
		return wrapErr(errTp, newWithLineError(runtime.CallerFileAndLine(1), fmt.Sprint(obj...)))
	case string:
		if strings.Contains(errTp, "%") {
			return newWithLineError(runtime.CallerFileAndLine(1), fmt.Sprintf(errTp, obj...))
		}
		if len(obj) == 0 {
			return newWithLineError(runtime.CallerFileAndLine(1), errTp)
		}
		return newWithLineError(runtime.CallerFileAndLine(1), fmt.Sprint(append([]interface{}{errTp + " "}, obj...)...))

	default:
		return newWithLineError(runtime.CallerFileAndLine(1), fmt.Sprint(append([]interface{}{errTp}, obj...)...))
	}

}

func wrapErr(errPre, errNow error) error {
	if errPre == nil {
		return nil
	}
	if errNow == nil {
		return nil
	}
	if _, ok := errNow.(WithLineError); !ok {
		errNow = newWithLineError(runtime.CallerFileAndLine(1), errNow)
	}
	if _, ok := errPre.(ErrorWrapper); !ok {
		errPre = newWithLineError(runtime.CallerFileAndLine(1), errPre)
	}
	return errPre.(ErrorWrapper).Wrap(errNow)
}
