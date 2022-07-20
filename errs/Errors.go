package errs

import (
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/runtime"
	"strings"
)

/*
WithLine Usage:
	1. WithLine(err) -> add caller line code to err
	2. WithLine(err1, err2) -> wrap err2 into err1 and add caller line code
	3. WithLine(err1, string1) -> make a withlineError of string1 and wrap it to err1
	4. WithLine(string1) -> make a withlineError of string1
	5. WithLine(string1 , obj ...) -> make a withlineError of fmt.Errorf(string1, obj...)
	6. WithLine(string1_with_no_"%" , obj ...) -> make a withlineError of fmt.Sprint(string1, obj...)
*/
func WithLine(err interface{}, obj ...interface{}) error {
	callerLine := runtime.CallerFileAndLine(1)
	errorCase := func(errTp error) error {
		if len(obj) == 0 {
			return errTp
			//resErr = errTp
			//break
		} else if len(obj) == 1 && obj[0] != nil {
			var obj0 error
			switch ob := obj[0].(type) {
			case error:
				obj0 = ob
			case string:
				obj0 = withLineError{lineCode: callerLine, error: errors.New(ob)}
			}
			return wrapError(errTp, obj0)
			//resErr = wrapError(errTp, obj0)
			//break
		}
		return wrapError(errTp, wrapError(errTp, withLineError{lineCode: callerLine, error: errors.New(fmt.Sprint(obj...))}))
	}
	var resErr error
	switch errTp := err.(type) {
	case withLineError:
		return withLineError{lineCode: callerLine, error: errorCase(errTp)}
	case error:
		resErr = errorCase(errTp)
		break
	case string:
		if len(obj) == 0 {
			resErr = errors.New(errTp)
			break
		} else if strings.Contains(errTp, "%") {
			resErr = fmt.Errorf(errTp, obj...)
			break
		} else {
			resErr = errors.New(fmt.Sprint(append([]interface{}{errTp + " "}, obj...)...))
			break
		}
	default:
		resErr = errors.New(fmt.Sprint(append([]interface{}{errTp}, obj...)...))
		break
	}
	return withLineError{lineCode: callerLine, error: resErr}
}

func wrapError(errParent error, errChild error) error {
	return fmt.Errorf("{ \n\t%w\n\t=> %s \n}", errParent, errChild.Error())
}
