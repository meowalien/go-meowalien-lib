package redis

import (
	"fmt"
	"testing"
)

func TestFormatKey(t *testing.T) {
	fmt.Println(KeyFormatter{
		Split:  DefaultSplit,
		Prefix: DefaultPrefix,
	}.Format("a", "b", "c"))
}
