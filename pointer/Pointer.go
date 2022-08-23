package pointer

func ToInterfacePointer[T any](i T) *T {
	return &i
}
