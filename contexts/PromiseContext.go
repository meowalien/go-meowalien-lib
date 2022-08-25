package contexts

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

// Make a new PromiseContext with the given name and WaitGroup wg,
// the wg could be nil, if so, the context will act as context.Context
func NewPromiseContext(parent PromiseContext, wg *sync.WaitGroup) (ctx PromiseContext, cancel context.CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := promiseContext{PromiseContext: parent, childWaitGroup: wg}
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, context.Canceled) }
}

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type DoneStd interface {
	doneStd() (chFc <-chan struct{})
}

type PromiseDone interface {
	Done() (chFc <-chan func())
}

// PromiseContext will add 1 to the WaitGroup when the Done Called, and minus 1 when the Done returned function called
type PromiseContext interface {
	Deadline() (deadline time.Time, ok bool)
	Err() error
	DoneStd
	PromiseDone
}

type canceler interface {
	cancel(removeFromParent bool, err error)
	PromiseDone
	DoneStd
}
type promiseContext struct {
	PromiseContext
	mu             sync.Mutex
	done           atomic.Value
	children       map[canceler]struct{}
	err            error
	childWaitGroup *sync.WaitGroup
}

func (c *promiseContext) doneStd() <-chan struct{} {
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

func (c *promiseContext) Done() (chFc <-chan func()) {
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

func (c *promiseContext) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

type stringer interface {
	String() string
}

func contextName(c PromiseContext) string {
	if s, ok := c.(stringer); ok {
		return s.String()
	}
	return reflect.TypeOf(c).String()
}

func (c *promiseContext) String() string {
	return contextName(c.PromiseContext) + ".WithCancel"
}

func (c *promiseContext) cancel(removeFromParent bool, err error) {
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
		removeChild(c.PromiseContext, c)
	}
}

var goroutines int32

func propagateCancel(parent PromiseContext, child canceler) {
	done := parent.Done()
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
			case okFc := <-parent.Done():
				child.cancel(false, parent.Err())
				okFc()
			case okFc1 := <-parent.Done():
				okFc1()
			}
		}()
	}
}

var cancelCtxKey int

func parentCancelCtx(parent PromiseContext) (*promiseContext, bool) {
	done := parent.doneStd()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.(*promiseContext)
	if !ok {
		return nil, false
	}
	pdone, _ := p.done.Load().(chan struct{})
	if pdone != done {
		return nil, false
	}
	return p, true
}

func removeChild(parent PromiseContext, child canceler) {
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

type emptyCtx int

func (e emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (e emptyCtx) Err() error {
	return nil
}

func (e emptyCtx) doneStd() (chFc <-chan struct{}) {
	return nil
}

func (e emptyCtx) Done() (chFc <-chan func()) {
	return nil
}
func (e *emptyCtx) String() string {
	return "emptyCtx"
}

var nilCtx = new(emptyCtx)

func NilCtx() *emptyCtx {
	return nilCtx
}
