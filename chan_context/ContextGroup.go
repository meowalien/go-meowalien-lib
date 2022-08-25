package chan_context

import (
	"context"
	"sync"
)

type ContextGroup interface {
	Close()
	NewContext() (WaitContext, context.CancelFunc)
	Child(name string) ContextGroup
}

func RootContextGroup(name string) ContextGroup {
	cg := newContextGroup(name, nilCtx)
	return cg
}

func newContextGroup(name string, ctx WaitContext) ContextGroup {
	cg := &contextGroup{
		name: name,
	}
	cg.ctx, cg.cancel = newWaitContext(ctx, nil)

	return cg
}

type contextGroup struct {
	wg     sync.WaitGroup
	name   string
	cancel context.CancelFunc
	ctx    WaitContext
	child  []ContextGroup
	lock   sync.Mutex
	closed bool
}

func (c *contextGroup) NewContext() (ctx WaitContext, cancel context.CancelFunc) {
	return newWaitContext(c.ctx, &c.wg)
}

func (c *contextGroup) Child(name string) ContextGroup {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		panic("context group closed")
	}

	newGroup := newContextGroup(name, c.ctx)

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
