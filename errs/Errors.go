package errs

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/runtime"
	"strings"
)

type withLineError struct {
	lineCode string
	error
}

func (w withLineError) Unwrap() error {
	return w.error
}

func (w withLineError) Error() string {
	return fmt.Sprintf("%s: %v", w.lineCode, w.error)
}

func addLineFormat(lineCode string, err interface{}) withLineError {
	switch errTp := err.(type) {
	case withLineError:
		return withLineError{lineCode: lineCode, error: errTp}
	case error:
		return withLineError{lineCode: lineCode, error: errTp}
	default:
		return withLineError{lineCode: lineCode, error: fmt.Errorf("%v", err)}
	}
}

func wrapErrorFormat(errParent error, errChild error) error {
	return fmt.Errorf("{ \n\t%w\n\t=> %s \n}", errParent, errChild.Error())
}

func WithLine(err interface{}, obj ...interface{}) error {
	callerLine := runtime.CallerFileAndLine(1)
	switch errTp := err.(type) {
	case error:
		if len(obj) == 0 {
			return addLineFormat(callerLine, errTp)
		} else if len(obj) == 1 && obj[0] != nil {
			var obj0 error
			switch ob := obj[0].(type) {
			case error:
				obj0 = ob
			case string:
				obj0 = addLineFormat(callerLine, ob)
			}
			return addLineFormat(callerLine, wrapErrorFormat(errTp, obj0))
		}
		return addLineFormat(callerLine, wrapErrorFormat(errTp, addLineFormat(callerLine, fmt.Sprint(obj...))))
	case string:
		if strings.Contains(errTp, "%") {
			return addLineFormat(callerLine, fmt.Sprintf(errTp, obj...))
		}
		if len(obj) == 0 {
			return addLineFormat(callerLine, errTp)
		}
		return addLineFormat(callerLine, fmt.Sprint(append([]interface{}{errTp + " "}, obj...)...))

	default:
		return addLineFormat(callerLine, fmt.Sprint(append([]interface{}{errTp}, obj...)...))
	}
}
