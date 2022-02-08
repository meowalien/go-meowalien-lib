package errs

import (
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/runtime"
	"go.uber.org/zap/buffer"
)

type WithLineError interface {
	error
	Wrap(err ...error) *withLineError
	Msg(msg ...interface{}) *withLineError
}

type withLineError struct {
	preErr *withLineError
	nowErr error
	msg    string
}

func (w *withLineError) Error() string {
	return w.GetChain().String()
}
func (w *withLineError) String() string {
	return w.Error()
}

func (w *withLineError) Msg(msg ...interface{}) *withLineError {
	if msg == nil || len(msg) == 0 {
		return w
	}

	if w.msg == "" {
		w.msg = fmt.Sprint(msg...)
	} else {
		w.msg = fmt.Sprintf("%s <- %s", w.msg, fmt.Sprint(msg...))
	}
	return w
}

func (w *withLineError) Wrap(err ...error) *withLineError {
	if err == nil || len(err) == 0 {
		return w
	}
	newE := w
	for _, err2 := range err {
		newE = &withLineError{
			preErr: w,
			nowErr: err2,
		}
	}
	return newE
}

var strBufferPool = buffer.NewPool()

func (w *withLineError) GetChain() *buffer.Buffer {
	var a *buffer.Buffer
	if w.preErr == nil {
		a = strBufferPool.Get()
	} else {
		a = w.preErr.GetChain()
		a.AppendString(" > ")
	}

	a.AppendString(w.nowErr.Error())
	if w.msg != "" {
		a.AppendString(fmt.Sprintf("( %s )", w.msg))
	}
	return a
}

func new(line string, err error) *withLineError {
	if line == "" {
		return &withLineError{
			preErr: nil,
			nowErr: err,
		}
	}
	return &withLineError{
		preErr: nil,
		nowErr: fmt.Errorf("#%s: %s", line, err.Error()),
	}
}

var Loop = false

func WithLine(err interface{}, msg ...interface{}) WithLineError {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *withLineError:
		if Loop {
			return new(runtime.CallerFileAndLine(1), e).Msg(msg...)
		} else {
			if msg == nil || len(msg) == 0 {
				return e
			} else {
				return e.Msg(msg...)
			}
		}
	case error:
		return new(runtime.CallerFileAndLine(1), e).Msg(msg...)
	case string:
		return new(runtime.CallerFileAndLine(1), errors.New(fmt.Sprintf(e, msg...)))
	default:
		return new(runtime.CallerFileAndLine(1), errors.New(fmt.Sprintf(fmt.Sprintf("%v", e), msg...)))
	}
}
