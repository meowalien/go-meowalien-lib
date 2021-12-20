package password

import (
	"testing"
)

//var bf = buffer.Buffer{}

func TestHashPassword(t *testing.T) {
	doit := func() string{
		bf := bffPool.Get()
		defer bf.Free()
		defer bf.Reset()
		bf.WriteString("1234")
		bf.WriteString("fkfkfkfk")
		return bf.String()
	}
	for i := 0; i < 10000; i++ {
		if doit() != "1234fkfkfkfk"{
			t.Error("fail on doit")
		}
	}
}

