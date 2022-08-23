package lock

import "sync"

type MapLocker[T comparable] interface {
	Lock(key T) func()
}

func NewMapLocker[T comparable]() MapLocker[T] {
	return &mapLocker[T]{}
}

var rwMutexPool = sync.Pool{New: func() interface{} {
	return &sync.RWMutex{}
}}

type mapLocker[T comparable] struct {
	sMap sync.Map
}

func (m *mapLocker[T]) Lock(key T) func() {
	lk := rwMutexPool.Get().(*sync.RWMutex)
	l, loaded := m.sMap.LoadOrStore(key, lk)
	if loaded {
		rwMutexPool.Put(lk)
	}
	l.(*sync.RWMutex).Lock()
	return l.(*sync.RWMutex).Unlock
}

func (m *mapLocker[T]) Free(key T) (loaded bool) {
	l, loaded := m.sMap.LoadAndDelete(key)
	if loaded {
		rwMutexPool.Put(l)
	}
	return false
}
