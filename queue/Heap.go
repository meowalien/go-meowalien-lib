package queue

type heap []SortableItem

func (h heap) Len() int { return len(h) }
func (h heap) Cap() int { return cap(h) }
func (h heap) Less(i, j int) bool {
	return h[i].Less(h[j])
}
func (h heap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *heap) Push(x interface{}) {
	*h = append(*h, x.(SortableItem))
}

func (h *heap) Peek() interface{} {
	return (*h)[0]
}
func (h *heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
