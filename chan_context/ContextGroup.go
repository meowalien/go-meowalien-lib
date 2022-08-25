package chan_context

import (
	"context"
	"sync"
)

type GroupContext[T comparable] interface {
	Key() T
	Close()
	Context() (WaitContext, context.CancelFunc)
	Child(name T) GroupContext[T]
}

func NewContextGroup[T comparable](name T) GroupContext[T] {
	return newGroupContext(name, nilCtx)
}

func newGroupContext[T comparable](name T, ctx WaitContext) GroupContext[T] {
	cg := &contextGroup[T]{
		name: name,
	}
	cg.ctx, cg.cancel = newWaitContext(ctx, nil)

	return cg
}

type contextGroup[T comparable] struct {
	wg     sync.WaitGroup
	name   T
	cancel context.CancelFunc
	ctx    WaitContext
	child  []GroupContext[T]
	lock   sync.Mutex
	closed bool
}

func (c *contextGroup[T]) Key() T {
	return c.name
}

func (c *contextGroup[T]) Context() (ctx WaitContext, cancel context.CancelFunc) {
	return newWaitContext(c.ctx, &c.wg)
}

func (c *contextGroup[T]) Child(name T) GroupContext[T] {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		panic("context group closed")
	}

	newGroup := newGroupContext(name, c.ctx)

	c.child = append(c.child, newGroup)
	return newGroup
}

func (c *contextGroup[T]) Close() {
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
