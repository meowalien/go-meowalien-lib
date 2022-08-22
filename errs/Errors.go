package errs

import (
	"github.com/meowalien/go-meowalien-lib/runtime"
)

/*
New Usage:
	New(any...) => make a withLineError of errors.New(fmt.Sprint(any...))
	New(string , any...) => make a withLineError of fmt.Errorf(string, obj...)
	New(error1 , error2/string , error3/string ...) => wrap error1(error2(error3 ...)))
	New(error) make a withLineError of error
*/
var New = func(err any, obj ...any) *withLineError {
	return newWithLineErrorFromAny(true, err, runtime.Caller(1), obj...)
}
