package chan_context

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type DoneStd interface {
	doneStd() (chFc <-chan struct{})
}
type DonePromise interface {
	DonePromise() (chFc <-chan struct{}, ok func())
}
type WaitContext interface {
	context.Context
	DoneStd
	DonePromise
}

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

var (
	//background = new(emptyCtx)
	todo = new(emptyCtx)
)

type CancelFunc func()

func WithCancel(parent WaitContext, wg *sync.WaitGroup) (ctx WaitContext, cancel CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := newCancelCtx(parent, wg)
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, context.Canceled) }
}

func newCancelCtx(parent WaitContext, wg *sync.WaitGroup) cancelCtx {
	return cancelCtx{WaitContext: parent, childWaitGroup: wg}
}

var goroutines int32

func propagateCancel(parent WaitContext, child canceler) {
	done, okFc := parent.DonePromise()
	if done == nil {
		return // parent is never canceled
	}
	defer okFc()

	select {
	case <-done:
		// parent is already canceled
		child.cancel(false, parent.Err())
		return
	default:
	}

	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil {
			// parent has already been canceled
			child.cancel(false, p.err)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		atomic.AddInt32(&goroutines, +1)
		go func() {
			doneC, okFc := parent.DonePromise()
			defer okFc()
			doneC1, okFc1 := parent.DonePromise()
			defer okFc1()
			select {
			case <-doneC:
				child.cancel(false, parent.Err())
				//defer doneOK()
			case <-doneC1:
				//defer doneOK()
			}
		}()
	}
}

var cancelCtxKey int

func parentCancelCtx(parent WaitContext) (*cancelCtx, bool) {
	done := parent.doneStd()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	pdone, _ := p.done.Load().(chan struct{})
	if pdone != done {
		return nil, false
	}
	return p, true
}

func removeChild(parent WaitContext, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}

type canceler interface {
	cancel(removeFromParent bool, err error)
	DonePromise
	DoneStd
}

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type cancelCtx struct {
	WaitContext

	mu             sync.Mutex
	done           atomic.Value
	children       map[canceler]struct{}
	err            error
	childWaitGroup *sync.WaitGroup
}

func (c *cancelCtx) Value(key any) any {
	if key == &cancelCtxKey {
		return c
	}
	return value(c.WaitContext, key)
}

func (c *cancelCtx) doneStd() <-chan struct{} {
	d := c.done.Load()
	if d != nil {
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	d = c.done.Load()
	if d == nil {
		d = make(chan struct{})
		c.done.Store(d)
	}
	return d.(chan struct{})
}

//var count int64

func (c *cancelCtx) DonePromise() (chFc <-chan struct{}, ok func()) {
	if c.childWaitGroup != nil {
		c.childWaitGroup.Add(1)
	}
	chFc = c.doneStd()
	ok = func() {
		if c.childWaitGroup != nil {
			c.childWaitGroup.Done()
		}
	}
	return
}
func (c *cancelCtx) Done() (ch <-chan struct{}) {
	ch, ok := c.DonePromise()
	ok()
	return
}

func (c *cancelCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

type stringer interface {
	String() string
}

func contextName(c WaitContext) string {
	if s, ok := c.(stringer); ok {
		return s.String()
	}
	return reflect.TypeOf(c).String()
}

func (c *cancelCtx) String() string {
	return contextName(c.WaitContext) + ".WithCancel"
}

func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}
	c.err = err
	d, _ := c.done.Load().(chan struct{})
	if d == nil {
		c.done.Store(closedchan)
	} else {
		close(d)
	}
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.WaitContext, c)
	}
}

type timerCtx struct {
	cancelCtx
	timer *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) String() string {
	return contextName(c.cancelCtx.WaitContext) + ".WithDeadline(" +
		c.deadline.String() + " [" +
		time.Until(c.deadline).String() + "])"
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		// Remove this timerCtx from its parent cancelCtx's children.
		removeChild(c.cancelCtx.WaitContext, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}

//func WithTimeout(parent Context, timeout time.Duration, wg *sync.WaitGroup) (Context, CancelFunc) {
//	return WithDeadline(parent, time.Now().Add(timeout), wg)
//}
//
//func WithValue(parent Context, key, val any) Context {
//	if parent == nil {
//		panic("cannot create context from nil parent")
//	}
//	if key == nil {
//		panic("nil key")
//	}
//	if !reflect.TypeOf(key).Comparable() {
//		panic("key is not comparable")
//	}
//	return &valueCtx{parent, key, val}
//}

type valueCtx struct {
	WaitContext
	key, val any
}

func stringify(v any) string {
	switch s := v.(type) {
	case stringer:
		return s.String()
	case string:
		return s
	}
	return "<not Stringer>"
}

func (c *valueCtx) String() string {
	return contextName(c.WaitContext) + ".WithValue(type " +
		reflect.TypeOf(c.key).String() +
		", val " + stringify(c.val) + ")"
}

func (c *valueCtx) Value(key any) any {
	if c.key == key {
		return c.val
	}
	return value(c.WaitContext, key)
}

func value(c WaitContext, key any) any {
	for {
		switch ctx := c.(type) {
		case *valueCtx:
			if key == ctx.key {
				return ctx.val
			}
			c = ctx.WaitContext
		case *cancelCtx:
			if key == &cancelCtxKey {
				return c
			}
			c = ctx.WaitContext
		case *timerCtx:
			if key == &cancelCtxKey {
				return &ctx.cancelCtx
			}
			c = ctx.WaitContext
		case *emptyCtx:
			return nil
		default:
			return c.Value(key)
		}
	}
}
