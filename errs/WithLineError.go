package errs

import (
	"fmt"
)

type ErrorWrapper interface {
	Wrap(err error) error
}

func newWithLineError(callerLocate string, err interface{}) WithLineError {
	return WithLineError{error: fmt.Errorf("%s: %v", callerLocate, err)}
}

type WithLineError struct {
	error
}

func (e WithLineError) Wrap(err error) error {
	return fmt.Errorf("{ %w } -> { %s }", e, err.Error())
}
