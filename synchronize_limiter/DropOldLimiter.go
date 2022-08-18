package synchronize_limiter

import (
	"context"
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync"
)

type dropOldLimiter struct {
	wait             sync.WaitGroup
	cancel           context.CancelFunc
	ctx              context.Context
	waitingTaskQueue chan func()
	runningThread    chan struct{}
	lock             sync.Mutex
}

func (s *dropOldLimiter) Do(ctx context.Context, f func(ctx context.Context)) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	select {
	case <-s.ctx.Done():
		err = errs.New("limiter stopping")
		return
	case s.waitingTaskQueue <- func() {
		f(ctx)
	}:
		fmt.Println("add to waiting queue")
		return
	case s.runningThread <- struct{}{}:
		fmt.Println("get thread")
		break
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			err = errs.New("timeout")
		} else {
			err = errs.New("limiter context done: %w", ctx)
		}
		return
	default:
		if cap(s.waitingTaskQueue) == 0 {
			err = errs.New("the WaitingQueueLimit should be greater than 0 in Strategy_DropOld")
			return
		}
		select {
		case <-s.waitingTaskQueue:
		default:
		}

		s.waitingTaskQueue <- func() {
			f(ctx)
		}
	}

	s.wait.Add(1)
	go func(f func(ctx context.Context)) {
		f(ctx)
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

func (s *dropOldLimiter) Stop() {
	s.cancel()
	s.Wait()
}

func (s *dropOldLimiter) Wait() {
	s.wait.Wait()
}
