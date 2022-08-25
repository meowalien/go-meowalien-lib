package slice

import "fmt"

func RemoveIdx[T any](s []T, n int) []T {
	return append(s[:n], s[n+1:]...)
}

func RemoveMatch[T comparable](s []T, target T) []T {
	fmt.Println("before remove", s)
	for i := 0; i < len(s); i++ {
		if s[i] == target {
			s = RemoveIdx(s, i)
			i--
		}
	}
	fmt.Println("after remove", s)

	return s
}

func ToAnySlice[T any](sl []T) (ans []any) {
	ans = make([]interface{}, len(sl))
	for i, v := range sl {
		ans[i] = v
	}
	return
}
