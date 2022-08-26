package contexts

import (
	"context"
	"sync"
)

/*
	ContextGroup is a group of PromiseContext, it can create child ContextGroup,
	all child ContextGroup will be closed before the parent.
*/
type ContextGroup interface {
	PromiseContext
	Close()
	ChildGroup() ContextGroup
}

func NewContextGroup(ctx PromiseContext) ContextGroup {
	cg := &contextGroup{}
	cg.ctx, cg.cancel = NewPromiseContext(ctx, &sync.WaitGroup{})
	return cg
}

type contextGroup struct {
	wg     sync.WaitGroup
	cancel context.CancelFunc
	ctx    PromiseContext
	child  []ContextGroup
	lock   sync.Mutex
	closed bool
}

func (c *contextGroup) PromiseDone() (chFc <-chan func()) {
	ctx, _ := NewPromiseContext(c.ctx, &c.wg)
	return ctx.PromiseDone()
}

func (c *contextGroup) ChildGroup() ContextGroup {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		panic("context group closed")
	}

	newGroup := NewContextGroup(c.ctx)

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
