package maps

import (
	"sync/atomic"
	"unsafe"
)

func newEntryFlowMap[T any](i T, expungedFlowMap unsafe.Pointer) *entryFlowMap[T] {
	return &entryFlowMap[T]{p: unsafe.Pointer(&i), expungedFlowMap: expungedFlowMap}
}

type entryFlowMap[T any] struct {
	expungedFlowMap unsafe.Pointer
	p               unsafe.Pointer // *interface{}
}

func (e *entryFlowMap[T]) tryStore(i *T) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == e.expungedFlowMap {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryFlowMap[T]) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, e.expungedFlowMap, nil)
}

func (e *entryFlowMap[T]) storeLocked(i *T) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (e *entryFlowMap[T]) tryLoadOrStore(i T) (actual T, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == e.expungedFlowMap {
		return actual, false, false
	}
	if p != nil {
		return *(*T)(p), true, true
	}

	// DeapCopy the interface after the first load to make this method more amenable
	// to escape analysis: if we hit the "load" path or the entry is expunged, we
	// shouldn't bother heap-allocating.
	ic := i
	for {
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == e.expungedFlowMap {
			return actual, false, false
		}
		if p != nil {
			return *(*T)(p), true, true
		}
	}
}

func (e *entryFlowMap[T]) delete() (value T, ok bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == e.expungedFlowMap {
			return value, false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return *(*T)(p), true
		}
	}
}
func (e *entryFlowMap[T]) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, e.expungedFlowMap) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == e.expungedFlowMap
}
func (e *entryFlowMap[T]) load() (value T, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == e.expungedFlowMap {
		return value, false
	}
	return *(*T)(p), true
}
