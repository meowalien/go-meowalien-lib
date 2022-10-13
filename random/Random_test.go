package random

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRandomString(t *testing.T) {
	s := RandomString(10)
	fmt.Println(s)
}

func TestFakeRand(t *testing.T) {
	s := rand.NewSource(928953616732700779)

	fmt.Println(s.Int63())
	fmt.Println(s.Int63())
	fmt.Println(s.Int63())
	fmt.Println(s.Int63())
	fmt.Println(s.Int63())
	fmt.Println(s.Int63())
	fmt.Println(s.Int63())
	fmt.Println(s.Int63())
	fmt.Println(s.Int63())
	fmt.Println(s.Int63())
}
