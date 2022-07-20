package errs

import "fmt"

type WithLineError interface {
	error
	Unwrap() error
}

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
