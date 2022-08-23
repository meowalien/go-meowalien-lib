package queue

import "golang.org/x/exp/constraints"

type heap[T constraints.Ordered] []T

func (h heap[T]) Len() int { return len(h) }
func (h heap[T]) Cap() int { return cap(h) }
func (h heap[T]) Less(i, j int) bool {
	return h[i] < h[j]
}
func (h heap[T]) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *heap[T]) Push(x interface{}) {
	*h = append(*h, x.(T))
}

func (h *heap[T]) Peek() interface{} {
	return (*h)[0]
}
func (h *heap[T]) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
