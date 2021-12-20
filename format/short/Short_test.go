package short

import (
	"math/rand"
	"testing"
)

func TestShort(t *testing.T) {
	for i := 0; i < 10000; i++ {
		rd := rand.Int63()
		hx := Base58Short(rd)
		ax := DeBase58Short(hx)
		if ax != rd {
			t.Errorf("error ax != rd: %d != %d" , ax , rd)
		}
	}
	for i := 0; i < 10000; i++ {
		rd := rand.Int63()
		hx := MaxShort(rd)
		ax := DeMaxShort(hx)
		if ax != rd {
			t.Errorf("error ax != rd: %d != %d" , ax , rd)
		}
	}
}
