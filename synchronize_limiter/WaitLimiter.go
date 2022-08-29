package synchronize_limiter

import (
	"context"
	"errors"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync"
)

type waitLimiter struct {
	wait             sync.WaitGroup
	waitingTaskQueue chan func()
	stopChan         chan struct{}
	threadCount      int
	ctx              context.Context
	cancel           context.CancelFunc
}

func (s *waitLimiter) Stop(ctx context.Context) {
	s.cancel()
	for len(s.waitingTaskQueue) != 0 {
		//	wait Consumer to consume all task
	}
	close(s.stopChan)
	s.wait.Wait() // wait for queue to be empty

}

func (s *waitLimiter) Do(ctx context.Context, f func()) (err error) {
	select {
	case <-s.ctx.Done():
		err = errs.New("limiter stopping")
		return
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			err = errs.New("timeout")
		} else {
			err = errs.New("limiter context done: ", ctx.Err())
		}
		return
	case s.waitingTaskQueue <- f:
		return
	}
}
func (s *waitLimiter) startConsumer() {
	s.wait.Add(s.threadCount)
	for i := 0; i < s.threadCount; i++ {
		go func(i int) {
			defer s.wait.Done()
			for {
				select {
				case <-s.stopChan:
					return
				case f := <-s.waitingTaskQueue:
					f()
					continue
				}
			}
		}(i)
	}
}
