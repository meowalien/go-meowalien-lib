package trace

import (
	"context"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/runtime"
)

var NewBlock = func(ctx context.Context, name string) func() {
	fmt.Printf("%s: StartBlock %s \n", runtime.Caller(1, runtime.CALLER_FORMAT_LONG), name)
	return func() {
		fmt.Printf("%s: EndBlock %s \n", runtime.Caller(1, runtime.CALLER_FORMAT_LONG), name)
	}
}
