package global_pool

import "sync"



var rWMutexPool = sync.Pool{
	New: func() interface{} {
		return &sync.RWMutex{}
	},
}
func GetRWMutex() *sync.RWMutex {
	return rWMutexPool.Get().(*sync.RWMutex)
}

func PutRWMutex(c *sync.RWMutex) {
	rWMutexPool.Put(c)
}
