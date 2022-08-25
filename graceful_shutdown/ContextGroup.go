package graceful_shutdown

import (
	"context"
	"sync"
)

func newContextGroup(ctx *promiseDone) *contextGroup {
	cg := &contextGroup{}
	cg.ctx, cg.cancel = newPromiseDone(ctx, &sync.WaitGroup{})

	return cg
}

type contextGroup struct {
	wg     sync.WaitGroup
	cancel context.CancelFunc
	ctx    *promiseDone
	child  []*contextGroup
	lock   sync.Mutex
	closed bool
}

func (c *contextGroup) Done() (chFc <-chan func()) {
	ctx, _ := newPromiseDone(c.ctx, &c.wg)
	return ctx.Done()
}

func (c *contextGroup) childGroup() *contextGroup {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		panic("context group closed")
	}

	newGroup := newContextGroup(c.ctx)

	c.child = append(c.child, newGroup)
	return newGroup
}

func (c *contextGroup) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.closed = true
	if c.child != nil {
		for i := len(c.child) - 1; i >= 0; i-- {
			c.child[i].Close()
		}
	}
	c.cancel()
	c.wg.Wait()
}
