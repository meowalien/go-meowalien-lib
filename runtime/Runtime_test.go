package runtime

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/debug"
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
	debug.DebugMode = true
	s := Caller(0)
	fmt.Println(s)
}
