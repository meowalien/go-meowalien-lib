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
	switch errTp := err.(type) {
	case error:
		return withLineError{lineCode: callerLine, error: errorCase(callerLine, errTp, obj)}
	case string:
		var resErr error
		if len(obj) == 0 {
			resErr = errors.New(errTp)
		} else if strings.Contains(errTp, "%") {
			resErr = fmt.Errorf(errTp, obj...)
		} else {
			resErr = errors.New(fmt.Sprint(append([]interface{}{errTp + " "}, obj...)...))
		}
		return withLineError{lineCode: callerLine, error: resErr}
	default:
		resErr := errors.New(fmt.Sprint(append([]interface{}{errTp}, obj...)...))
		return withLineError{lineCode: callerLine, error: resErr}
	}
}

func errorCase(callerLine string, errTp error, obj []interface{}) error {
	if len(obj) == 0 {
		return errTp
	} else if len(obj) == 1 && obj[0] != nil {
		var obj0 error
		switch ob := obj[0].(type) {
		case error:
			obj0 = ob
		case string:
			obj0 = withLineError{lineCode: callerLine, error: errors.New(ob)}
		}
		return wrapError(errTp, obj0)
	}
	return wrapError(errTp, wrapError(errTp, withLineError{lineCode: callerLine, error: errors.New(fmt.Sprint(obj...))}))
}

func wrapError(errParent error, errChild error) error {
	return fmt.Errorf("{ \n\t%w\n\t=> %s \n}", errParent, errChild.Error())
}
