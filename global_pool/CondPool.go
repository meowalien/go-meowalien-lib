package global_pool

import "sync"

var rCondPool = sync.Pool{
	New: func() interface{} {
		return sync.NewCond(nil)
	},
}

func GetCond(mu *sync.RWMutex) *sync.Cond {
	cd := rCondPool.Get().(*sync.Cond)
	cd.L = mu
	return cd
}

func PutCond(c *sync.Cond) {
	rCondPool.Put(c)
}
