package queue

import "golang.org/x/exp/constraints"
import (
	sheep "container/heap"
)

type PriorityQueue[T constraints.Ordered] interface {
	Queue[T]
}

func NewPriorityQueue[T constraints.Ordered](cap int) PriorityQueue[T] {
	return &priorityQueue[T]{
		heap: make([]T, 0, cap),
	}
}

type priorityQueue[T constraints.Ordered] struct {
	heap[T]
}

func (p *priorityQueue[T]) Push(item T) {
	sheep.Push(&p.heap, item)
}

func (p *priorityQueue[T]) Pop() T {
	return sheep.Pop(&p.heap).(T)
}

func (p *priorityQueue[T]) Peek() T {
	return p.heap.Peek().(T)
}

func (p *priorityQueue[T]) Len() int {
	return p.heap.Len()
}
func (p *priorityQueue[T]) Cap() int {
	return p.heap.Cap()
}
