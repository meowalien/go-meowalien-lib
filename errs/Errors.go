package errs

import (
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/runtime"
	"strings"
)

/*
WithLine Usage:
	1. WithLine(err) -> add caller line code to err, if err is nil, return nil
	2. WithLine(err1, err2) -> wrap err2 into err1 and add caller line code
	3. WithLine(err1, string1) -> make a withlineError of string1 and wrap it to err1
	4. WithLine(string1) -> make a withlineError of string1
	5. WithLine(string1 , obj ...) -> make a withlineError of fmt.Errorf(string1, obj...)
	6. WithLine(string1_with_no_"%" , obj ...) -> make a withlineError of fmt.Sprint(string1, obj...)
*/
func WithLine(err interface{}, obj ...interface{}) error {
	if err == nil {
		return nil
	}
	callerLine := runtime.CallerFileAndLine(1)
	var resErr error
	switch errTp := err.(type) {
	case error:
		switch {
		case len(obj) == 0:
			resErr = errTp
		case len(obj) == 1 && obj[0] != nil:
			var obj0 error
			switch ob := obj[0].(type) {
			case error:
				obj0 = ob
			case string:
				obj0 = withLineError{lineCode: callerLine, error: errors.New(ob)}
			}
			resErr = wrapError(errTp, obj0)
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
