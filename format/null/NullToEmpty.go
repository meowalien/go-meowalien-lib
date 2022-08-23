package null

func Null[T any](s *T) (st T) {
	if s != nil {
		return *s
	}
	return
}
