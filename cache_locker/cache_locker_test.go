package cache_locker

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheLocker(t *testing.T) {
	var x int
	ch := NewCache(&x)
	wg := sync.WaitGroup{}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t, release := ch.Load()
			*t++
			time.Sleep(time.Millisecond / 10)
			release()
		}()
	}

	wg.Add(1)
	time.AfterFunc(time.Millisecond/5, func() {
		count := ch.UserCount()
		fmt.Println("count:", count)
		defer wg.Done()
		t, release := ch.Freeze()
		fmt.Println(*t)
		defer release()
	})
	wg.Wait()
}
