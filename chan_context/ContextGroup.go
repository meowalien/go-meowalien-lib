package chan_context

import (
	"fmt"
	"sync"
)

type ContextGroup interface {
	Close()
	NewContext() (Context, CancelFunc)
	Child(name string) ContextGroup
}

func RootContextGroup(name string) ContextGroup {
	cg := newContextGroup(name, todo)
	return cg
}

func newContextGroup(name string, ctx Context) ContextGroup {
	cg := &contextGroup{
		name: name,
	}
	cg.ctx, cg.cancel = WithCancel(ctx, nil)

	return cg
}

type contextGroup struct {
	wg     *sync.WaitGroup
	name   string
	cancel CancelFunc
	ctx    Context
	child  []ContextGroup
	lock   sync.Mutex
	closed bool
}

func (c *contextGroup) NewContext() (ctx Context, cancel CancelFunc) {
	ctx, cancel = WithCancel(c.ctx, c.wg)
	//wctx = newWaitContext(ctx, c.wg)
	return
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
	if c.wg != nil {
		c.wg.Wait()
	}

	fmt.Println("End Close: ", c.name)
	fmt.Println("----------------------------------------------------")
	//time.Sleep(time.Second * 2)
}
