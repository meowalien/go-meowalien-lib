package pointer

func ToInterfacePointer(i interface{}) *interface{} {
	return &i
}
