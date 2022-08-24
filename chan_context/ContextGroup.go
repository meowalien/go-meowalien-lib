package chan_context

import (
	"context"
	"fmt"
	"sync"
)

type ContextGroup interface {
	Close()
	NewContext() (WaitContext, context.CancelFunc)
	Child(name string) ContextGroup
}

func RootContextGroup(name string) ContextGroup {
	cg := newContextGroup(name, todo)
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
	fmt.Println("Start Close: ", c.name)
	c.closed = true
	if c.child != nil {
		for i := len(c.child) - 1; i >= 0; i-- {
			c.child[i].Close()
		}
		fmt.Println("closed all child: ", c.name)
	}
	c.cancel()
	//if c.wg != nil {
	c.wg.Wait()
	//}

	fmt.Println("End Close: ", c.name)
	fmt.Println("----------------------------------------------------")
	//time.Sleep(time.Second * 2)
}
