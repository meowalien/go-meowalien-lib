package errs

import "fmt"

//type WithLineError interface {
//	error
//	Unwrap() error
//}

type withLineError struct {
	lineCode string
	error
}

func (w withLineError) Unwrap() error {
	return w.error
}

func (w withLineError) Error() string {
	//return fmt.Sprintf("%s: %s", w.lineCode, tp.Error())
	switch tp := w.error.(type) { //nolint:errorlint
	case withLineError:
		return fmt.Sprintf("%s: \n\t%s", w.lineCode, tp.Error())
	default:
		return fmt.Sprintf("%s: %s", w.lineCode, tp.Error())
	}
}
