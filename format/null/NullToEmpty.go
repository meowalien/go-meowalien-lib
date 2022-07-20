package null

func NullString(s *string) (st string) {
	if s != nil {
		return *s
	}
	return
}
