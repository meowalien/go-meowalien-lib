package chan_context

import "time"

type emptyCtx int

func (e *emptyCtx) Done() <-chan struct{} {
	return nil
}

func (e *emptyCtx) DonePromise() (chFc <-chan struct{}, ok func()) {
	return nil, ok
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
	//case background:
	//	return "context.Background"
	case todo:
		return "context.TODO"
	default:
		panic("unreachable")
	}
	return "unknown empty Context"
}
