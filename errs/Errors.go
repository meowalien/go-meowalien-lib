package errs

import (
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/runtime"
	"strings"
)

/*
New Usage:
	1. New(err) -> add caller line code to err, if err is nil, return nil
	2. New(err1, err2) -> wrap err2 into err1 and add caller line code
	3. New(err1, string1) -> make a withlineError of string1 and wrap it to err1
	4. New(string1) -> make a withlineError of string1
	5. New(string1 , obj ...) -> make a withlineError of fmt.Errorf(string1, obj...)
	6. New(string1_with_no_"%" , obj ...) -> make a withlineError of fmt.Sprint(string1, obj...)
*/
func New(err interface{}, obj ...interface{}) error {
	if err == nil {
		if len(obj) != 1 {
			if len(obj) > 1 {
				panic("New: obj must be one or zero")
			}
			return nil
		}
		if e, ok := obj[0].(error); ok {
			err = e
			obj = obj[1:]
		}
	}
	callerLine := runtime.Caller(1)
	var resErr error
	switch errTp := err.(type) {
	case error:
		switch {
		case len(obj) == 0:
			resErr = errTp
		case len(obj) == 1:
			if obj[0] == nil {
				if errors.Is(errTp, withLineErrorType) {
					return errTp
				} else {
					resErr = errTp
					break
				}
			} else {
				var obj0 error
				switch ob := obj[0].(type) {
				case error:
					obj0 = ob
				case string:
					obj0 = withLineError{lineCode: callerLine, error: errors.New(ob)}
				}
				resErr = wrapError(errTp, obj0)
			}
		default:
			resErr = wrapError(errTp, wrapError(errTp, withLineError{lineCode: callerLine, error: errors.New(fmt.Sprint(obj...))}))
		}
	case string:
		switch {
		case len(obj) == 0:
			resErr = errors.New(errTp)
		case strings.Contains(errTp, "%"):
			resErr = fmt.Errorf(errTp, obj...)
		default:
			resErr = errors.New(fmt.Sprint(append([]interface{}{errTp + " "}, obj...)...))
		}
	default:
		resErr = errors.New(fmt.Sprint(append([]interface{}{errTp}, obj...)...))
	}
	return withLineError{lineCode: callerLine, error: resErr}
}

func wrapError(errParent error, errChild error) error {
	return fmt.Errorf("{ \n\t%w\n\t=> %s \n}", errParent, errChild.Error())
}
