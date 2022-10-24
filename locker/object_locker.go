package locker

import (
	"sync"
	"sync/atomic"
)

type ObjectLocker[T any] interface {
	Load() (t *T, release func())
	Store(T)
	Do(func(t *T))
	Freeze() (t *T, release func())
	UserCount() uint64
}

func NewObjectLocker[T any](t *T) ObjectLocker[T] {
	return &cacheLocker[T]{t: t}
}

type cacheLocker[T any] struct {
	t          *T
	lock       sync.RWMutex
	tUserCount uint64
}

func (c *cacheLocker[T]) UserCount() uint64 {
	return c.tUserCount
}
func (c *cacheLocker[T]) Store(t T) {
	c.lock.Lock()
	defer c.lock.Unlock()
	*c.t = t
}
func (c *cacheLocker[T]) Load() (t *T, release func()) {
	c.lock.RLock()
	atomic.AddUint64(&c.tUserCount, 1)
	return c.t, func() {
		atomic.AddUint64(&c.tUserCount, ^uint64(0))
		c.lock.RUnlock()
	}
}

func (c *cacheLocker[T]) Do(f func(t *T)) {
	t, release := c.Load()
	defer release()
	f(t)
}

func (c *cacheLocker[T]) Freeze() (t *T, release func()) {
	c.lock.Lock()
	return c.t, c.lock.Unlock
}
