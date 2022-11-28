package errs

import (
	"gitlab.geax.io/demeter/backendmodules/runtime"
)

/*
New Usage:

	New(any...) => make a withLineError of errors.New(fmt.Sprint(any...))
	New(string , any...) => make a withLineError of fmt.Errorf(string, obj...)
	New(error1(could be nil) , error2/string , error3/string ...) => wrap error1(error2(error3 ...)))
	New(error) make a withLineError of error
*/
var New = func(err any, obj ...any) WithLineError {
	ee := newWithLineErrorFromAny(true, err, runtime.Caller(1, runtime.CALLER_FORMAT_SHORT), obj...)
	// to make sure that the returned error is nil type and nil value
	if ee == nil {
		return nil
	}
	return ee
}
