package queue

import "golang.org/x/exp/constraints"

type Queue[T constraints.Ordered] interface {
	Push(T)
	Pop() T
	Peek() T
	Len() int
	Cap() int
}
