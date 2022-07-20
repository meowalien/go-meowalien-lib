package insertable_linked_list

import (
	"container/list"
	"github.com/meowalien/go-meowalien-lib/global_pool"
	"sync"
	//"sync/atomic"
)

var listPool = sync.Pool{New: func() interface{} {
	return &orderList{list: global_pool.GetList(), mu: global_pool.GetRWMutex()}
}}

func New() OrderList {
	return listPool.Get().(*orderList) //&orderList{list: global_pool.GetList(), mu: global_pool.GetRWMutex()}
}

type OrderList interface {
	Put(order int64, value interface{})
	Front() OrderElement
	GetAndRemoveLowerThen(cursor int64) (all []OrderElement)
	GetAllAndRemove() (all []OrderElement)
	Free()
}

type orderList struct {
	mu   *sync.RWMutex
	list *list.List
}

func (l *orderList) Free() {
	listPool.Put(l)
}

func (l *orderList) Put(order int64, value interface{}) {
	ele := NewElement(order, value)
	for checkEle := l.Front(); checkEle != nil; checkEle = checkEle.Next() {
		valueOrder := checkEle.Order()
		if valueOrder <= ele.Order() {
			continue
		} else {
			l.InsertBefore(ele, checkEle)
			return
		}
	}
	l.PushBack(ele)
}

func (l *orderList) GetAndRemoveLowerThen(cursor int64) (all []OrderElement) {
	for checkEle := l.Front(); checkEle != nil; {
		valueOrder := checkEle.Order()
		if valueOrder > cursor {
			return
		}

		nextCheckEle := checkEle.Next()
		all = append(all, checkEle)
		l.Remove(checkEle)
		checkEle = nextCheckEle
		continue

	}
	return
}

func (l *orderList) GetAllAndRemove() (all []OrderElement) {
	for checkEle := l.Front(); checkEle != nil; {
		nextCheckEle := checkEle.Next()
		all = append(all, checkEle)
		l.Remove(checkEle)
		checkEle = nextCheckEle
		continue
	}
	return
}

func (l *orderList) Front() OrderElement {
	f := l.list.Front()
	if f == nil {
		return nil
	}
	return f.Value.(OrderElement)
}

func (l *orderList) InsertBefore(ele OrderElement, ele2 OrderElement) OrderElement {
	l.mu.Lock()
	defer l.mu.Unlock()
	listElement := l.list.InsertBefore(ele, ele2.GetListElement())
	ele.SetListElement(listElement)
	return ele
}

func (l *orderList) PushBack(ele OrderElement) OrderElement {
	l.mu.Lock()
	defer l.mu.Unlock()
	listElement := l.list.PushBack(ele)
	ele.SetListElement(listElement)
	return ele
}

func (l *orderList) Remove(ele OrderElement) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.list.Remove(ele.GetListElement())
}
