package maps

type Map[K comparable, T any] interface {
	LoadAndDelete(key K) (value T, loaded bool)
	LoadOrStore(key K, value T) (actual T, loaded bool)
	Load(key K) (value T, ok bool)
	Store(key K, value T)
	Delete(key K)
	Range(f func(key K, value T) bool)
}
