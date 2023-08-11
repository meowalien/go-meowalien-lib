package hash

import (
	"testing"
)

type A struct {
	AA string
	BB int
}

func TestFastHash(t *testing.T) {
	h1, err := FastHash(A{
		AA: "Hello",
		BB: 123,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	h2, err := FastHash(A{
		AA: "Hello",
		BB: 123,
	})

	if err != nil {
		t.Fatal(err)
		return
	}
	if h1 != h2 {
		t.Fatal("h1 != h2")
		return
	}
}

//func TestFastHash2(t *testing.T) {
//	h1, err := FastHash2(A{
//		AA: "Hello",
//		BB: 123,
//	})
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	fmt.Println(h1)
//	h2, err := FastHash2(BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB{
//		AA: "Hello",
//		BB: 123,
//	})
//
//	if err != nil {
//		t.Fatal(err)
//		return
//	}
//	fmt.Println(h2)
//}
//
//type MyStruct struct {
//	Name  string
//	Value int
//}
//
////go test -bench . -benchtime=1000000x
//
//// Benchmark for FastHash
//func BenchmarkFastHash(b *testing.B) {
//	testStruct := MyStruct{Name: "Test", Value: 42}
//	b.ResetTimer() // Reset the timer so that the setup doesn't affect the results
//
//	for i := 0; i < b.N; i++ {
//		_, _ = FastHash(testStruct)
//	}
//}
//
//// Benchmark for FastHash1
//func BenchmarkFastHash1(b *testing.B) {
//	testStruct := MyStruct{Name: "Test", Value: 42}
//	b.ResetTimer() // Reset the timer so that the setup doesn't affect the results
//
//	for i := 0; i < b.N; i++ {
//		_, _ = FastHash1(testStruct)
//	}
//}
