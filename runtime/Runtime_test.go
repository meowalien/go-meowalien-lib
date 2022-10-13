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
	s := Caller(0, CALLER_FORMAT_SHORT)
	fmt.Println(s)
}

func TestShort(t *testing.T) {
	s := Caller(0, CALLER_FORMAT_SHORT)
	fmt.Println(s)
}

func TestLong(t *testing.T) {
	s := Caller(0, CALLER_FORMAT_LONG)
	fmt.Println(s)
}
