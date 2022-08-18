package synchronize_limiter

import (
	"context"
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync"
)

type waitLimiter struct {
	wait             sync.WaitGroup
	waitingTaskQueue chan func()
	runningThread    chan struct{}
	ctx              context.Context
	cancel           context.CancelFunc
}

func (s *waitLimiter) Stop() {
	s.cancel()
	s.Wait()
}

func (s *waitLimiter) Wait() {
	s.wait.Wait()
}

func (s *waitLimiter) Do(ctx context.Context, f func()) (err error) {
	select {
	case <-s.ctx.Done():
		err = errs.New("limiter stopping")
		return
	case s.waitingTaskQueue <- f:
		fmt.Println("add to waiting queue")
		return
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			err = errs.New("timeout")
		} else {
			err = errs.New("limiter context done: %w", ctx)
		}
		return
	case s.runningThread <- struct{}{}:
		fmt.Println("get thread")
		break
	}

	s.wait.Add(1)
	go func(f func()) {
		f()
		for {
			select {
			case <-s.ctx.Done():
				err = errs.New("limiter stopping")
				return
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					err = errs.New("timeout")
				} else {
					err = errs.New("limiter context done: %w", ctx)
				}
				return
			case nextf := <-s.waitingTaskQueue:
				fmt.Println("run from waiting queue")
				nextf()
				continue
			default:
				fmt.Println("release running thread")
				<-s.runningThread
				s.wait.Done()
				return
			}
		}
	}(f)
	return
}
