package wrapper

type Wrapper[T any] interface {
	Wrap(cli T) T
}

func Wrap[T any](cli T, wps ...Wrapper[T]) T {
	for _, wp := range wps {
		cli = wp.Wrap(cli)
	}
	return cli
}
