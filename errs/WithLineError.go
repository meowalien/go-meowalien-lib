package errs

import (
	"errors"
	"fmt"
	"strings"
)

type withLineError struct {
	parent   *withLineError
	lineCode string
	error
	layer int
}

func (w *withLineError) Unwrap() error {
	if w.parent == nil {
		return w.error
	} else {
		return w.parent
	}
}

func (w *withLineError) Error() (s string) {
	tabs := strings.Repeat("\t", w.Layer())
	s = fmt.Sprintf("%s%s: %s", tabs, w.lineCode, w.error.Error())
	if w.parent == nil {
		return
	} else {
		return w.parent.FormatChild(s)
	}

}

func (w withLineError) Layer() (l int) {
	return w.layer
}

func (w *withLineError) Wrap(a any, caller string) (res *withLineError) {
	ne := newWithLineErrorFromAny(false, a, caller)
	ne.layer = w.layer + 1
	ne.parent = w
	return ne
}

func (w *withLineError) FormatChild(childErrStr string) (res string) {
	tabs := strings.Repeat("\t", w.Layer())
	res = fmt.Sprintf("%s%s: %s {\n%s%s\n%s}",
		tabs,
		w.lineCode,
		w.error,
		tabs,
		childErrStr,
		tabs,
	)
	if w.parent == nil {
		return
	} else {
		return w.parent.FormatChild(res)
	}
}

func (w withLineError) Deliver(caller string) *withLineError {
	w.lineCode = fmt.Sprintf("%s => %s", w.lineCode, caller)
	x := w
	return &x
}

func newWithLineErrorFromAny(deliver bool, err any, caller string, obj ...any) *withLineError {
	if err == nil {
		if len(obj) == 0 {
			return nil
		}
		return newWithLineErrorFromAny(deliver, obj[0], caller, obj[1:]...)
	}
	switch errTp := err.(type) {
	case string:
		return newWithLineErrorFromError(fmt.Errorf(errTp, obj...), caller)
	case error:
		var parentErr *withLineError
		parentErr, ok := errTp.(*withLineError) //nolint:errorlint
		if !ok {
			parentErr = newWithLineErrorFromError(errTp, caller)
		} else if deliver {
			parentErr = parentErr.Deliver(caller)
		}

		for _, a := range obj {
			if a == nil {
				continue
			}
			parentErr = parentErr.Wrap(a, caller)
		}
		return parentErr
	default:
		return newWithLineErrorFromError(errors.New(fmt.Sprint(append([]any{errTp}, obj...)...)), caller)
	}
}

func newWithLineErrorFromError(err error, caller string) *withLineError {
	return &withLineError{lineCode: caller, error: err}
}
