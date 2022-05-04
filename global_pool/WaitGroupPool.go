package global_pool

import "sync"

var waitGroupPool = sync.Pool{
	New: func() interface{} {
		return &sync.WaitGroup{}
	},
}
func GetWaitGroup() *sync.WaitGroup {
	return waitGroupPool.Get().(*sync.WaitGroup)
}

func PutWaitGroup(c *sync.WaitGroup) {
	waitGroupPool.Put(c)
}
