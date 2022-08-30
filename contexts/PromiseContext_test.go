package contexts

import (
	"context"
	"testing"
	"time"
)

func TestPromiseContext(t *testing.T) {
	//promiseCtx, promiseCancel := NewPromiseContext(nil, &sync.WaitGroup{})
	ctx, cancel := context.WithCancel(context.Background())
	promiseCtx := NewContextGroup(nil)
	go func() {
		<-ctx.Done()
		promiseCtx.Close()
	}()

	readPumpCtx := NewContextGroup(promiseCtx)
	writePumpCtx := NewContextGroup(promiseCtx)
	go func() {
		ok := <-promiseCtx.PromiseDone()
		t.Log("promiseCtx done")
		ok()
	}()

	go func() {
		ok := <-readPumpCtx.PromiseDone()
		t.Log("readPumpCtx done")
		ok()
	}()
	go func() {
		ok := <-writePumpCtx.PromiseDone()
		t.Log("writePumpCtx done")
		ok()
	}()

	time.Sleep(time.Second * 1)
	cancel()
}
