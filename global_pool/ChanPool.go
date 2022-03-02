package global_pool

import "sync"

var int64ChanPool = sync.Pool{
	New: func() interface{} {
		return make(chan int64, 500)
	},
}

func GetInt64ChanPool() chan int64 {
	return int64ChanPool.Get().(chan int64)
}

func PutInt64ChanPool(c chan int64) {
	int64ChanPool.Put(c)
}
