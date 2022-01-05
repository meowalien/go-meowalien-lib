package short

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestShort(t *testing.T) {
	for i := 1; i < 10000; i++ {
		//rd := rand.Int63()
		hx := Base58Short(int64(i))

		ax := DeBase58Short(hx)
		fmt.Println("raw: ",i)
		fmt.Println("ax: ",ax)
		if ax != int64(i) {
			t.Errorf("error ax != rd: %d != %d" , ax , i)
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
