package convert

import "unsafe"

func AtoB[T any, B any](a T) (b B) {
	return *(*B)(unsafe.Pointer(&a)) //nolimt:gosec
}
