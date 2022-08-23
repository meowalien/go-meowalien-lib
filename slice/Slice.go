package slice

func Remove[T any](s []T, n int) []T {
	return append(s[:n], s[n+1:]...)
}

func ToAnySlice[T any](sl []T) (ans []any) {
	ans = make([]interface{}, len(sl))
	for i, v := range sl {
		ans[i] = v
	}
	return
}
