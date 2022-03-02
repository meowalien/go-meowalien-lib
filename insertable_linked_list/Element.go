package insertable_linked_list


type Element struct {
	Order int64
	Value interface{}
}

//type Element interface {
//	Next() Element
//	Prev() Element
//	PutOnOrder(order int64, v interface{})
//	GetBeforeIncludeOrder(order int64, sum []interface{})
//	isRoot() bool
//	iterate(c chan<- interface{})
//	pushOrder()
//	pullOrder()
//	Value() interface{}
//}
//
//// Element is an element of a linked list.
//type element struct {
//	order      int64
//	next, prev *element
//	list       *List
//	value      interface{}
//}
//
//func (e *element) Value() interface{} {
//	return e.value
//}
//
//// Next returns the next list element or nil.
//func (e *element) Next() Element {
//	if p := e.next; e.list != nil && p != &e.list.root {
//		return p
//	}
//	return nil
//}
//
//func (e *element) Prev() Element {
//	if p := e.prev; e.list != nil && p != &e.list.root {
//		return p
//	}
//	return nil
//}
//
//func (e *element) PutOnOrder(targetOrder int64, v interface{}) {
//	fmt.Printf("in PutOnOrder current : %d , target : %d \n", e.order, targetOrder)
//	if e.isRoot() && e.next == e {
//		fmt.Println("is root")
//		e.list.InsertAfter(v, e)
//		return
//	}
//
//	if e.order < targetOrder {
//		if !e.next.isRoot() {
//			e.next.PutOnOrder(targetOrder, v)
//			return
//		}else{
//		//	對列最後一個並且order還是小於 targetOrder
//			e.list.InsertAfterNotPush(v, e)
//			return
//		}
//	} else if e.order == targetOrder {
//		fmt.Println("e.order == targetOrder")
//
//		e.list.InsertAfter(v, e)
//		return
//	}else if e.order > targetOrder{
//		//	最後一個的order大於指定order
//		fmt.Println("e.order > targetOrder")
//		e.list.InsertAfterNotPush(v, e)
//		return
//	}
//
//	//if e.isRoot() {
//	//	fmt.Println("aaa")
//	//	e.list.PushFront(v)
//	//	return
//	//}
//	//if e.next.isRoot() || e.order == order {
//	//	fmt.Println("bbb")
//	//	e.list.InsertAfter(v, e)
//	//	return
//	//}
//
//	//if e.order < order {
//	//	e.next.PutOnOrder(order, v)
//	//	return
//	//} else {
//	//	e.list.InsertBefore(v, e)
//	//	return
//	//}
//}
//
//func (e *element) GetBeforeIncludeOrder(order int64, sum []interface{}) {
//	if e.order <= order {
//		sum = append(sum, e.value)
//		if !e.next.isRoot() {
//			e.next.GetBeforeIncludeOrder(order, sum)
//		}
//		return
//	}
//}
//
//func (e *element) isRoot() bool {
//	//fmt.Println("e.list: ",e.list)
//	return e == &e.list.root
//}
//
//func (e *element) iterate(c chan<- interface{}) {
//	if !e.isRoot() {
//		fmt.Println("iterate-e order:", e.order)
//		c <- e.value
//	}
//
//	if e.next.isRoot() {
//		close(c)
//	} else {
//		e.next.iterate(c)
//	}
//
//}
//
//func (e *element) pushOrder() {
//	if e == nil {
//		fmt.Println("e == nil")
//		return
//	}
//	e.order += 1
//	fmt.Println("pushOrder: ", e.order)
//	if !e.next.isRoot() {
//		e.next.pushOrder()
//	}
//}
//
//func (e *element) pullOrder() {
//	e.order -= 1
//	//e.order -= 1
//	if !e.next.isRoot() {
//		e.next.pushOrder()
//	}
//}
