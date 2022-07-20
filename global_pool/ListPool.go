package global_pool

import (
	"container/list"
	"sync"
)

var listPool = sync.Pool{New: func() interface{} {
	return list.New()
}}

func GetList() *list.List {
	return listPool.Get().(*list.List)
}

func PutList(b *list.List) {
	listPool.Put(b)
}
