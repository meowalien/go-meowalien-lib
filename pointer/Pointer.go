package pointer

func ToInterfacePointer(i interface{}) *interface{}{
	x:= interface{}(i)
	return &x
}