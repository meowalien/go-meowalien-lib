package hash

import (
	"fmt"
	"testing"
)

type A struct {
	AA string
	BB int
}

type B struct {
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
	fmt.Println(h1)
	h2, err := FastHash(A{
		AA: "Hello",
		BB: 123,
	})

	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(h2)
}
