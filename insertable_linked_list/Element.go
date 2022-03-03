package insertable_linked_list

import "container/list"

func NewElement(order int64, value interface{}) *orderElement {
	return &orderElement{
		order: order,
		value: value,
	}
}

type OrderElement interface {
	Order() int64
	Value() interface{}
	Next() OrderElement
	SetListElement(element *list.Element)
	GetListElement() (element *list.Element)
}

type orderElement struct {
	order       int64
	value       interface{}
	listElement *list.Element
}

func (e *orderElement) SetListElement(element *list.Element) {
	e.listElement = element
}

func (e *orderElement) GetListElement() (element *list.Element) {
	return e.listElement
}

func (e *orderElement) Next() OrderElement {
	f :=  e.listElement.Next()
	if f == nil{
		return nil
	}
	return f.Value.(OrderElement)
}

func (e *orderElement) Order() int64 {
	return e.order
}

func (e *orderElement) Value() interface{} {
	return e.value
}
