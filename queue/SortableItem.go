package queue

type SortableItem interface {
	Less(other SortableItem) bool
}

type Indexable interface {
	Index() int
}

type IntSortableItem int

func (i IntSortableItem) Index() int {
	return int(i)
}

// Less if other less than self
func (i IntSortableItem) Less(other Indexable) bool {
	return int(i) > other.Index()
}
