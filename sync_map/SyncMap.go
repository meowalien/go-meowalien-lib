package sync_map

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type SyncMap[K comparable, T any] interface {
	Load(key K) (value T, ok bool)
	Store(key K, value T)
	LoadOrStore(key K, value T) (actual T, loaded bool)
	LoadAndDelete(key K) (value T, loaded bool)
	Delete(key K)
	Range(f func(key K, value T) bool)
}

func NewMap[K comparable, T any]() SyncMap[K, T] {
	return &syncMap[K, T]{
		expungedFlowMap: unsafe.Pointer(new(T)),
	}
}

type syncMap[K comparable, T any] struct {
	expungedFlowMap unsafe.Pointer
	mu              sync.Mutex
	read            atomic.Value // readOnly
	dirty           map[K]*entryFlowMap[T]
	misses          int
}

type readOnlyFlowMap[K comparable, T any] struct {
	m       map[K]*entryFlowMap[T]
	amended bool // true if the dirty map contains some key not in m.
}

func (m *syncMap[K, T]) Load(key K) (value T, ok bool) {
	read, _ := m.read.Load().(readOnlyFlowMap[K, T])
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		// Avoid reporting a spurious miss if m.dirty got promoted while we were
		// blocked on m.mu. (If further loads of the same key will not miss, it's
		// not worth copying the dirty map for this key.)
		read, _ = m.read.Load().(readOnlyFlowMap[K, T])
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// Regardless of whether the entry was present, record a miss: this key
			// will take the slow path until the dirty map is promoted to the read
			// map.
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return value, false
	}
	return e.load()
}

func (m *syncMap[K, T]) Store(key K, value T) {
	read, _ := m.read.Load().(readOnlyFlowMap[K, T])
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyFlowMap[K, T])
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			// The entry was previously expunged, which implies that there is a
			// non-nil dirty map and this entry is not in it.
			m.dirty[key] = e
		}
		e.storeLocked(&value)
	} else if e, ok := m.dirty[key]; ok {
		e.storeLocked(&value)
	} else {
		if !read.amended {
			// We're adding the first new key to the dirty map.
			// Make sure it is allocated and mark the read-only map as incomplete.
			m.dirtyLocked()
			m.read.Store(readOnlyFlowMap[K, T]{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryFlowMap(value, m.expungedFlowMap)
	}
	m.mu.Unlock()
}

func (m *syncMap[K, T]) LoadOrStore(key K, value T) (actual T, loaded bool) {
	// Avoid locking if it's a clean hit.
	read, _ := m.read.Load().(readOnlyFlowMap[K, T])
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyFlowMap[K, T])
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok {
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	} else {
		if !read.amended {
			// We're adding the first new key to the dirty map.
			// Make sure it is allocated and mark the read-only map as incomplete.
			m.dirtyLocked()
			m.read.Store(readOnlyFlowMap[K, T]{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryFlowMap(value, m.expungedFlowMap)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

func (m *syncMap[K, T]) LoadAndDelete(key K) (value T, loaded bool) {
	read, _ := m.read.Load().(readOnlyFlowMap[K, T])
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyFlowMap[K, T])
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			delete(m.dirty, key)
			// Regardless of whether the entry was present, record a miss: this key
			// will take the slow path until the dirty map is promoted to the read
			// map.
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if ok {
		return e.delete()
	}
	return value, false
}

func (m *syncMap[K, T]) Delete(key K) {
	m.LoadAndDelete(key)
}

func (m *syncMap[K, T]) Range(f func(key K, value T) bool) {
	// We need to be able to iterate over all of the keys that were already
	// present at the start of the call to Range.
	// If read.amended is false, then read.m satisfies that property without
	// requiring us to hold m.mu for a long time.
	read, _ := m.read.Load().(readOnlyFlowMap[K, T])
	if read.amended {
		// m.dirty contains keys not in read.m. Fortunately, Range is already O(N)
		// (assuming the caller does not break out early), so a call to Range
		// amortizes an entire copy of the map: we can promote the dirty copy
		// immediately!
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyFlowMap[K, T])
		if read.amended {
			read = readOnlyFlowMap[K, T]{m: m.dirty}
			m.read.Store(read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

func (m *syncMap[K, T]) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyFlowMap[K, T]{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *syncMap[K, T]) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyFlowMap[K, T])
	m.dirty = make(map[K]*entryFlowMap[T], len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}
