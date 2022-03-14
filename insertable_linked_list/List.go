package insertable_linked_list

import (
	"container/list"
	"sync"
	//"sync/atomic"
)

func New() OrderList {
	return &orderList{list: new(list.List).Init()}
}

type OrderList interface {
	Put(order int64, value interface{})
	Front() OrderElement
	GetAndRemoveLowerThen(cursor int64) (all []OrderElement)
	GetAllAndRemove() (all []OrderElement)
}

type orderList struct {
	mu   sync.Mutex
	list *list.List
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
		all = append(all, checkEle.(OrderElement))
		l.Remove(checkEle)
		checkEle = nextCheckEle
		continue

	}
	return
}

func (l *orderList) GetAllAndRemove() (all []OrderElement) {
	for checkEle := l.Front(); checkEle != nil; {
		nextCheckEle := checkEle.Next()
		all = append(all, checkEle.(OrderElement))
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
