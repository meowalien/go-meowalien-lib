package chan_context

import "time"

type emptyCtx int

var nilCtx = new(emptyCtx)

func (e *emptyCtx) Done() <-chan struct{} {
	return nil
}

func (e *emptyCtx) DonePromise() (chFc <-chan func()) {
	return nil
}

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (*emptyCtx) doneStd() <-chan struct{} {
	return nil
}

func (*emptyCtx) Err() error {
	return nil
}

func (*emptyCtx) Value(key any) any {
	return nil
}

func (e *emptyCtx) String() string {
	switch e {
	case nilCtx:
		return "context.TODO"
	default:
		panic("unreachable")
	}
}
