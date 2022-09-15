package thread

import (
	"context"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync/atomic"
	"testing"
	"time"
)

func TestDropOldestLimiter(t *testing.T) {

	round := 99999
	capacity := round / 10
	limiter := NewSynchronizeLimiter(Config{
		QueueFullStrategy:    Strategy_DropOldest,
		WaitingQueueCapacity: capacity,
		RunningThreadAmount:  1,
	})
	var count int64 = 0

	for i := 0; i < round; i++ {
		//ii := i
		time.Sleep(time.Microsecond * 1)
		err := limiter.Do(context.Background(), func() {
			time.Sleep(time.Microsecond * 10)
			//fmt.Println("AAAAAAAAAAAAAAAAAAA: ", ii)
			atomic.AddInt64(&count, 1)
		})
		if err != nil {
			panic(errs.New(err))
			return
		}
	}
	time.Sleep(time.Second * 1)
	limiter.WaitQueueClean()
	fmt.Println("count: ", count)
}
