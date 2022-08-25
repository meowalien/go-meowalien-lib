package contexts

import (
	"context"
	"sync"
)

// ContextGroup will promise to cancel all child context before parent context cancel.
type ContextGroup[T comparable] interface {
	Key() T
	Close()
	Child(name T) ContextGroup[T]
	PromiseDone
}

func NewContextGroup[T comparable](name T) ContextGroup[T] {
	return newContextGroup(name, NilCtx())
}

func newContextGroup[T comparable](name T, ctx PromiseContext) ContextGroup[T] {
	cg := &contextGroup[T]{
		name: name,
	}
	cg.ctx, cg.cancel = NewPromiseContext(ctx, &sync.WaitGroup{})

	return cg
}

type contextGroup[T comparable] struct {
	wg     sync.WaitGroup
	name   T
	cancel context.CancelFunc
	ctx    PromiseContext
	child  []ContextGroup[T]
	lock   sync.Mutex
	closed bool
}

func (c *contextGroup[T]) Key() T {
	return c.name
}
func (c *contextGroup[T]) Done() (chFc <-chan func()) {
	ctx, _ := NewPromiseContext(c.ctx, &c.wg)
	return ctx.Done()
}

//func (c *contextGroup[T]) Context() (ctx PromiseContext) {
//	ctx, _ = NewPromiseContext(c.ctx, &c.wg)
//	return
//}

func (c *contextGroup[T]) Child(name T) ContextGroup[T] {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		panic("context group closed")
	}

	newGroup := newContextGroup(name, c.ctx)

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
