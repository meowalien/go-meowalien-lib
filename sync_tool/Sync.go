package sync_tool

import "sync"

type MapLocker struct {
	sMap sync.Map
}

var rwMutexPool = sync.Pool{New: func() interface{} {
	return &sync.RWMutex{}
}}

func (m *MapLocker) Get(key interface{}) (lk *sync.RWMutex, loaded bool) {
	lk = rwMutexPool.Get().(*sync.RWMutex)
	l, loaded := m.sMap.LoadOrStore(key, lk)
	if loaded {
		rwMutexPool.Put(lk)
		lk = l.(*sync.RWMutex)
	}
	return
}

func (m *MapLocker) Free(key interface{}) (loaded bool) {
	l, loaded := m.sMap.LoadAndDelete(key)
	if loaded {
		rwMutexPool.Put(l)
	}
	return false
}

var maplocker = MapLocker{}

func Get(key interface{}) (lk *sync.RWMutex, exist bool) {
	return maplocker.Get(key)
}

func Free(key interface{}) (exist bool) {
	return maplocker.Free(key)
}
