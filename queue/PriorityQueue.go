package queue

import (
	sheep "container/heap"
)

type PriorityQueue struct {
	heap
}

func (p *PriorityQueue) Push(item SortableItem) {
	sheep.Push(&p.heap, item)
}

func (p *PriorityQueue) Pop() SortableItem {
	return sheep.Pop(&p.heap).(SortableItem)
}

func (p *PriorityQueue) Peek() SortableItem {
	return p.heap.Peek().(SortableItem)
}

func (p *PriorityQueue) Len() int {
	return p.heap.Len()
}
func (p *PriorityQueue) Cap() int {
	return p.heap.Cap()
}
