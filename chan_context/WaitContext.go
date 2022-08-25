package chan_context

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
)

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type DoneStd interface {
	doneStd() (chFc <-chan struct{})
}

type DonePromise interface {
	DonePromise() (chFc <-chan func())
}

type WaitContext interface {
	context.Context
	DoneStd
	DonePromise
}

type canceler interface {
	cancel(removeFromParent bool, err error)
	DonePromise
	DoneStd
}
type waitContext struct {
	WaitContext
	mu             sync.Mutex
	done           atomic.Value
	children       map[canceler]struct{}
	err            error
	childWaitGroup *sync.WaitGroup
}

func (c *waitContext) Value(key any) any {
	if key == &cancelCtxKey {
		return c
	}
	return value(c.WaitContext, key)
}

func (c *waitContext) doneStd() <-chan struct{} {
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

func (c *waitContext) DonePromise() (chFc <-chan func()) {
	nChFc := make(chan func())
	chFc = nChFc
	if c.childWaitGroup != nil {
		c.childWaitGroup.Add(1)
	}
	ch := c.doneStd()

	go func(ch <-chan struct{}) {
		<-ch
		select {
		case nChFc <- func() {
			if c.childWaitGroup != nil {
				c.childWaitGroup.Done()
			}
		}:
		default:
			if c.childWaitGroup != nil {
				c.childWaitGroup.Done()
			}
		}
	}(ch)
	return
}

func (c *waitContext) Done() (ch <-chan struct{}) {
	f := <-c.DonePromise()
	ch = closedchan
	f()
	return
}

func (c *waitContext) Err() error {
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

func (c *waitContext) String() string {
	return contextName(c.WaitContext) + ".WithCancel"
}

func (c *waitContext) cancel(removeFromParent bool, err error) {
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

func newWaitContext(parent WaitContext, wg *sync.WaitGroup) (ctx WaitContext, cancel context.CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := waitContext{WaitContext: parent, childWaitGroup: wg}
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, context.Canceled) }
}

var goroutines int32

func propagateCancel(parent WaitContext, child canceler) {
	done := parent.DonePromise()
	if done == nil {
		return // parent is never canceled
	}

	select {
	case okFc := <-done:
		// parent is already canceled
		child.cancel(false, parent.Err())
		okFc()
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
			select {
			case okFc := <-parent.DonePromise():
				child.cancel(false, parent.Err())
				okFc()
			case okFc1 := <-parent.DonePromise():
				okFc1()
			}
		}()
	}
}

var cancelCtxKey int

func parentCancelCtx(parent WaitContext) (*waitContext, bool) {
	done := parent.doneStd()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*waitContext)
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
