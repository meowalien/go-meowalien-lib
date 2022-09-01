package contexts

import (
	"context"
	"fmt"
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

func TestWithValue(t *testing.T) {
	readPumpCtx := NewContextGroup(nil)
	val := context.WithValue(readPumpCtx, "key1", "value1")
	val = context.WithValue(val, "key2", "value2")
	x := val.Value("key2")
	fmt.Println(x)
}

func TestWithTimeout(t *testing.T) {
	readPumpCtx := NewContextGroup(nil)
	val, _ := context.WithTimeout(readPumpCtx, time.Hour*200)
	val, _ = context.WithTimeout(val, time.Hour*200)
	d, x := val.Deadline()
	fmt.Println(d)
	fmt.Println(x)
}
