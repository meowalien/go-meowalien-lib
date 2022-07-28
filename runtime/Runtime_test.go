package runtime

import (
	"fmt"
	"testing"
)

func TestCallerFileAndLine(t *testing.T) {
	s := CallerStackTrace(0)
	fmt.Println(s)
}
func TestCaller(t *testing.T) {
	s := Caller(0)
	fmt.Println(s)
}

func TestCallerWhenAwaysStackTrace(t *testing.T) {
	AlwaysStackTrace = true
	s := Caller(0)
	fmt.Println(s)
}
