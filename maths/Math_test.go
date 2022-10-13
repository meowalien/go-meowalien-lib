package maths

import (
	"fmt"
	"testing"
)

func TestRound(t *testing.T) {
	number := 12.3456789

	fmt.Println(Round(number, 2))
	fmt.Println(Round(number, 3))
	fmt.Println(Round(number, 4))
	fmt.Println(Round(number, 5))

	number = -12.3456789
	fmt.Println(Round(number, 0))
	fmt.Println(Round(number, 1))
	fmt.Println(Round(number, 10))
}
