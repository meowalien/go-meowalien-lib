package bitmask

import (
	"fmt"
	"testing"
)

const (
	Bitmask_A Bitmask = 1 << iota
	Bitmask_B
	Bitmask_C
	Bitmask_A_B = Bitmask_A | Bitmask_B
)

func TestBitmask(t *testing.T) {
	fmt.Println(Bitmask_A_B.Has(Bitmask_A))
	fmt.Println(Bitmask_A_B.Has(Bitmask_B))
	fmt.Println(Bitmask_A.Has(Bitmask_B))
	fmt.Println(Bitmask_A_B.Has(Bitmask_C))
}
