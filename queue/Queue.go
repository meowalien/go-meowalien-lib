package queue

type Queue interface {
	Push(SortableItem)
	Pop() SortableItem
	Peek() SortableItem
	Len() int
	Cap() int
}
