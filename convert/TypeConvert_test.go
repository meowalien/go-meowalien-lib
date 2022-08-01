package convert

import (
	"testing"
	"unsafe"
)

func BenchmarkAtoB(b *testing.B) {
	//BenchmarkAtoB-4   	1000000000	         0.2984 ns/op
	for i := 0; i < b.N; i++ {
		AtoB[string, []byte]("abc")
	}
}
func BenchmarkAtoB_old(b *testing.B) {
	//BenchmarkAtoB_old-4   	1000000000	         0.3042 ns/op
	for i := 0; i < b.N; i++ {
		wildPosStr := "abc"
		_ = *(*[]byte)(unsafe.Pointer(&wildPosStr))
	}
}
