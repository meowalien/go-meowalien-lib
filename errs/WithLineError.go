package errs

import (
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/bitmask"
	"strings"
)

func newWithLineErrorFromAny(deliver bool, err any, caller string, obj ...any) *withLineError {
	if err == nil || err == (*withLineError)(nil) {
		if len(obj) == 0 {
			return nil
		} else if len(obj) == 1 {
			return newWithLineErrorFromAny(deliver, obj[0], caller)
		} else {
			return newWithLineErrorFromAny(deliver, obj[0], caller, obj[1:]...)
		}
	}
	switch errTp := err.(type) {
	case string:
		if strings.Contains(errTp, "%") {
			return newWithLineErrorFromError(fmt.Errorf(errTp, obj...), caller)
		} else {
			return newWithLineErrorFromError(errors.New(fmt.Sprint(append([]any{errTp}, obj...)...)), caller)
		}
	case error:
		var parentErr *withLineError
		parentErr, ok := errTp.(*withLineError) //nolint:errorlint
		if !ok {
			parentErr = newWithLineErrorFromError(errTp, caller)
		} else if deliver {
			parentErr = parentErr.deliver(caller)
		}

		for _, a := range obj {
			if a == nil {
				continue
			}
			parentErr = parentErr.wrap(a, caller)
		}
		return parentErr
	default:
		return newWithLineErrorFromError(errors.New(fmt.Sprint(append([]any{errTp}, obj...)...)), caller)
	}
}

func newWithLineErrorFromError(err error, caller string) *withLineError {
	return &withLineError{caller: caller, error: err}
}

type WithLineError interface {
	error
	WithCode(code bitmask.Bitmask) WithLineError
	HasCode(code bitmask.Bitmask) bool
	ErrorCode() bitmask.Bitmask
	Is(e error) bool
}

type withLineError struct {
	error
	parent    *withLineError
	caller    string
	layer     int
	errorCode bitmask.Bitmask
}

func (w *withLineError) ErrorCode() bitmask.Bitmask {
	return w.errorCode
}

func (w *withLineError) HasCode(code bitmask.Bitmask) bool {
	if w.errorCode == nil || code == nil {
		return false
	}
	return w.errorCode.Has(code)
}

func (w withLineError) WithCode(code bitmask.Bitmask) WithLineError {
	if code == nil {
		panic("using nil or 0 as error code")
	}
	//fmt.Println("w.errorCode : ", w.errorCode)
	//fmt.Println("code : ", code)
	if w.errorCode == nil {
		w.errorCode = code
	} else {
		w.errorCode = w.errorCode.Add(code)
	}
	return &w
}

func (w *withLineError) Is(e error) bool {
	if w.error == e {
		return true
	}
	ee, ok := e.(WithLineError)
	if ok {
		return ee.HasCode(w.errorCode)
	}
	if w.parent == nil {
		return false
	} else {
		return w.parent.Is(e)
	}
}

func (w *withLineError) Unwrap() error {
	if w.parent == nil {
		return w.error
	} else {
		return w.parent
	}
}

func (w *withLineError) Error() (s string) {
	tabs := strings.Repeat("\t", w.layer)
	s = fmt.Sprintf("%s%s: %s", tabs, w.caller, w.error.Error())
	if w.parent == nil {
		return
	} else {
		return w.parent.formatChild(s)
	}
}

func (w *withLineError) wrap(a any, caller string) (res *withLineError) {
	ne := newWithLineErrorFromAny(false, a, caller)
	ne.layer = w.layer + 1
	ne.parent = w
	return ne
}
func (w *withLineError) formatChild(childErrStr string) (res string) {
	tabs := strings.Repeat("\t", w.layer)
	res = fmt.Sprintf("%s%s: %s {\n%s%s\n%s}",
		tabs,
		w.caller,
		w.error,
		tabs,
		childErrStr,
		tabs,
	)
	if w.parent == nil {
		return
	} else {
		return w.parent.formatChild(res)
	}
}

func (w withLineError) deliver(caller string) *withLineError {
	w.caller = fmt.Sprintf("%s => %s", w.caller, caller)
	x := w
	return &x
}
