package chan_context

import (
	"time"
)

type emptyCtx int

func (e emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (e emptyCtx) Done() <-chan struct{} {
	return nil
}

func (e emptyCtx) Err() error {
	return nil
}

func (e emptyCtx) Value(key any) any {
	return nil
}

func (e emptyCtx) doneStd() (chFc <-chan struct{}) {
	return nil
}

func (e emptyCtx) DonePromise() (chFc <-chan func()) {
	return nil
}
func (e *emptyCtx) String() string {
	return "emptyCtx"
}

var nilCtx = new(emptyCtx)
