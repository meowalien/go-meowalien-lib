package maps

import "sync"

type LockMap[K comparable, T any] interface {
	Map[K, T]
	Len() int
	RLock()
	TryRLock() bool
	RUnlock()
	Lock()
	TryLock() bool
	Unlock()
}

func NewLockMap[K comparable, T any]() LockMap[K, T] {
	return &lockedMap[K, T]{
		m: make(map[K]T),
	}
}

type lockedMap[K comparable, T any] struct {
	sync.RWMutex
	m map[K]T
}

func (l *lockedMap[K, T]) Len() int {
	return len(l.m)
}

func (l *lockedMap[K, T]) LoadAndDelete(key K) (value T, loaded bool) {
	l.Lock()
	defer l.Unlock()
	value, loaded = l.m[key]
	if loaded {
		delete(l.m, key)
	}
	return
}

func (l *lockedMap[K, T]) LoadOrStore(key K, value T) (actual T, loaded bool) {
	l.Lock()
	defer l.Unlock()
	actual, loaded = l.m[key]
	if !loaded {
		l.m[key] = value
		actual = value
	}
	return
}

func (l *lockedMap[K, T]) Load(key K) (value T, ok bool) {
	l.RLock()
	defer l.RUnlock()
	value, ok = l.m[key]
	return
}

func (l *lockedMap[K, T]) Store(key K, value T) {
	l.Lock()
	defer l.Unlock()
	l.m[key] = value
}

func (l *lockedMap[K, T]) Delete(key K) {
	l.Lock()
	defer l.Unlock()
	delete(l.m, key)
}

func (l *lockedMap[K, T]) Range(f func(key K, value T) bool) {
	l.RLock()
	defer l.RUnlock()
	for key, value := range l.m {
		if !f(key, value) {
			return
		}
	}
}
