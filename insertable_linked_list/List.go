package insertable_linked_list

import (
	"container/list"
	"sync"
	//"sync/atomic"
)

func New() InsertableLinkedList {
	return new(List).Init()
}

type InsertableLinkedList interface {
	Put(cursor int64, c  *Element)
	Iterator() chan  *Element
	GetAndRemoveLowerThen(cursor int64) [] *Element
}


type List struct {
	mu sync.Mutex
	*list.List
}

func (l *List) Put(order int64,ele *Element) {
	if  l.Len() == 0 {
		l.PushBack(ele)
	}else{
		checkEle := l.Front()
		for {

			valueOrder := checkEle.Value.(*Element).Order
			if valueOrder <= ele.Order{
				checkEleNext := checkEle.Next()
				if checkEleNext != nil{
					checkEle = checkEleNext
					continue
				}else{
					l.InsertAfter(ele,checkEle)
					return
				}
			}else if valueOrder == ele.Order{
				l.InsertAfter(ele,checkEle)
				return
			}else if valueOrder > ele.Order{
				l.InsertBefore(ele,checkEle)
				return
			}

		}
	}
}

func (l *List) Iterator() chan  *Element {
	panic("implement me")
}

func (l *List) GetAndRemoveLowerThen(cursor int64) [] *Element {
	panic("implement me")
}

func (l *List) Init() InsertableLinkedList {
	return &List{List:l.List.Init()}
}

//
//func (l *List) Put(order int64, c interface{}) {
//	fmt.Println("put")
//	l.root.PutOnOrder(order, c)
//}
//
//func (l *List) Iterator() chan interface{} {
//	c := make(chan interface{}, l.Len())
//	l.root.iterate(c)
//	return c
//}
//
//func (l *List) GetAndRemoveLowerThen(order int64) []interface{} {
//	lst := make([]interface{}, 0, 20)
//	l.root.GetBeforeIncludeOrder(order, lst)
//	return lst
//}
//
//// Init initializes or clears list l.
//func (l *List) Init() *List {
//	l.mu = sync.Mutex{}
//	l.root.list = l
//	l.root.next = &l.root
//	l.root.prev = &l.root
//	l.len = 0
//	return l
//}
//
//// New returns an initialized list.
////func New() *List { return new(List).Init() }
//
//// Len returns the number of elements of list l.
//// The complexity is O(1).
//func (l *List) Len() int { return l.len }
//
//// Front returns the first element of list l or nil if the list is empty.
//func (l *List) Front() Element {
//	if l.len == 0 {
//		return nil
//	}
//	return l.root.next
//}
//
//// Back returns the last element of list l or nil if the list is empty.
//func (l *List) Back() Element {
//	if l.len == 0 {
//		return nil
//	}
//	return l.root.prev
//}
//
//// lazyInit lazily initializes a zero List value.
//func (l *List) lazyInit() {
//	if l.root.next == nil {
//		l.Init()
//	}
//}
//
//// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
//func (l *List) insertValue(v interface{}, at *element) *element {
//	return l.insert(&element{value: v}, at)
//}
//
//// insert inserts e after at, increments l.len, and returns e.
//func (l *List) insert(e, at *element) *element {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//
//	e.order = at.order
//
//	e.prev = at
//	e.next = at.next
//
//	e.prev.next = e
//	e.next.prev = e
//
//	e.list = l
//	l.len++
//
//	return e
//}
//
//// remove removes e from its list, decrements l.len, and returns e.
//func (l *List) remove(e *element) *element {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	e.next.pullOrder()
//	e.prev.next = e.next
//	e.next.prev = e.prev
//	e.next = nil // avoid memory leaks
//	e.prev = nil // avoid memory leaks
//	e.list = nil
//	l.len--
//	return e
//}
//
//// move moves e to next to at and returns e.
//func (l *List) move(e, at *element) *element {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//
//	if e == at {
//		return e
//	}
//	e.next.pullOrder()
//	e.prev.next = e.next
//	e.next.prev = e.prev
//
//	at.order = at.next.order
//	at.next.pushOrder()
//
//	e.prev = at
//	e.next = at.next
//	e.prev.next = e
//	e.next.prev = e
//
//	return e
//}
//
//// Remove removes e from l if e is an element of list l.
//// It returns the element value e.Value.
//// The element must not be nil.
//func (l *List) Remove(e *element) interface{} {
//	if e.list == l {
//		// if e.list == l, l must have been initialized when e was inserted
//		// in l or l == nil (e is a zero Element) and l.remove will crash
//		l.remove(e)
//	}
//	return e.value
//}
//
//// PushFront inserts a new element e with value v at the front of list l and returns e.
//func (l *List) PushFront(v interface{}) *element {
//	l.lazyInit()
//	return l.insertValue(v, &l.root)
//}
//
//// PushBack inserts a new element e with value v at the back of list l and returns e.
//func (l *List) PushBack(v interface{}) *element {
//	l.lazyInit()
//	return l.insertValue(v, l.root.prev)
//}
//
//// InsertBefore inserts a new element e with value v immediately before mark and returns e.
//// If mark is not an element of l, the list is not modified.
//// The mark must not be nil.
//func (l *List) InsertBefore(v interface{}, mark *element) *element {
//	if mark.list != l {
//		return nil
//	}
//	// see comment in List.Remove about initialization of l
//	return l.insertValue(v, mark.prev)
//}
//
//func (l *List) InsertAfterNotPush(v interface{}, mark *element) *element {
//	if mark.list != l {
//		return nil
//	}
//	// see comment in List.Remove about initialization of l
//	return l.insertValue(v, mark)
//}
//
//func (l *List) InsertAfter(v interface{}, mark *element) *element {
//	newElement := l.InsertAfterNotPush(v, mark)
//
//	if newElement.next.isRoot() {
//		newElement.pushOrder()
//	} else {
//		newElement.next.pushOrder()
//	}
//	return newElement
//}
//
//// MoveToFront moves element e to the front of list l.
//// If e is not an element of l, the list is not modified.
//// The element must not be nil.
//func (l *List) MoveToFront(e *element) {
//	if e.list != l || l.root.next == e {
//		return
//	}
//	// see comment in List.Remove about initialization of l
//	l.move(e, &l.root)
//}
//
//// MoveToBack moves element e to the back of list l.
//// If e is not an element of l, the list is not modified.
//// The element must not be nil.
//func (l *List) MoveToBack(e *element) {
//	if e.list != l || l.root.prev == e {
//		return
//	}
//	// see comment in List.Remove about initialization of l
//	l.move(e, l.root.prev)
//}
//
//// MoveBefore moves element e to its new position before mark.
//// If e or mark is not an element of l, or e == mark, the list is not modified.
//// The element and mark must not be nil.
//func (l *List) MoveBefore(e, mark *element) {
//	if e.list != l || e == mark || mark.list != l {
//		return
//	}
//	l.move(e, mark.prev)
//}
//
//// MoveAfter moves element e to its new position after mark.
//// If e or mark is not an element of l, or e == mark, the list is not modified.
//// The element and mark must not be nil.
//func (l *List) MoveAfter(e, mark *element) {
//	if e.list != l || e == mark || mark.list != l {
//		return
//	}
//	l.move(e, mark)
//}
//
//// PushBackList inserts a copy of another list at the back of list l.
//// The lists l and other may be the same. They must not be nil.
//func (l *List) PushBackList(other *List) {
//	l.lazyInit()
//	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
//		l.insertValue(e.Value(), l.root.prev)
//	}
//}
//
//// PushFrontList inserts a copy of another list at the front of list l.
//// The lists l and other may be the same. They must not be nil.
//func (l *List) PushFrontList(other *List) {
//	l.lazyInit()
//	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
//		l.insertValue(e.Value(), &l.root)
//	}
//}
