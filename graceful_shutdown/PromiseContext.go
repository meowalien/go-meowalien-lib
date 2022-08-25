package graceful_shutdown

import (
	"context"
	"github.com/meowalien/go-meowalien-lib/slice"
	"sync"
	"sync/atomic"
)

// Make a new PromiseContext with the given name and WaitGroup wg,
// the wg could be nil, if so, the context will act as context.Context
func newPromiseDone(parent *promiseDone, wg *sync.WaitGroup) (ctx *promiseDone, cancel context.CancelFunc) {
	c := promiseDone{promiseDone: parent, childWaitGroup: wg}
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, context.Canceled) }
}

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

// PromiseContext will add 1 to the WaitGroup when the Done Called,
// and minus 1 when the Done returned function called.
type promiseDone struct {
	*promiseDone
	mu             sync.Mutex
	doneVal        atomic.Value
	children       []*promiseDone //map[*promiseDone]struct{}
	err            error
	childWaitGroup *sync.WaitGroup
}

func (c *promiseDone) done() <-chan struct{} {
	d := c.doneVal.Load()
	if d != nil {
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	d = c.doneVal.Load()
	if d == nil {
		d = make(chan struct{})
		c.doneVal.Store(d)
	}
	return d.(chan struct{})
}

func (c *promiseDone) Done() (chFc <-chan func()) {
	nChFc := make(chan func())
	chFc = nChFc
	if c.childWaitGroup != nil {
		c.childWaitGroup.Add(1)
	}
	ch := c.done()

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

func (c *promiseDone) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *promiseDone) String() string {
	if c.promiseDone != nil {
		return c.promiseDone.String() + ".WithCancel"
	}
	return "emptyCtx" + ".WithCancel"
}

func (c *promiseDone) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}
	c.err = err
	d, _ := c.doneVal.Load().(chan struct{})
	if d == nil {
		c.doneVal.Store(closedchan)
	} else {
		close(d)
	}
	for i := len(c.children) - 1; i >= 0; i-- {
		// NOTE: acquiring the child's lock while holding parent's lock.
		c.children[i].cancel(false, err)
	}

	//for child := range c.children {
	//	// NOTE: acquiring the child's lock while holding parent's lock.
	//	child.cancel(false, err)
	//}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.promiseDone, c)
	}
}

var goroutines int32

func propagateCancel(parent *promiseDone, child *promiseDone) {
	if parent == nil {
		return
	}
	done := parent.Done()

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
			//if p.children == nil {
			//	p.children = make(map[*promiseDone]struct{})
			//}
			p.children = append(p.children, child)
			//p.children[child] = struct{}{}
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

func parentCancelCtx(parent *promiseDone) (*promiseDone, bool) {
	if parent == nil {
		return nil, false
	}
	done := parent.done()
	if done == closedchan || done == nil {
		return nil, false
	}

	pdone, _ := parent.doneVal.Load().(chan struct{})
	if pdone != done {
		return nil, false
	}
	return parent, true
}

func removeChild(parent *promiseDone, child *promiseDone) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		slice.RemoveMatch(p.children, child)
		//delete(p.children, child)
	}
	p.mu.Unlock()
}
