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
	s.cleanup()
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
		fmt.Println("add to waiting queue")
		return

	case s.runningThread <- struct{}{}:
		fmt.Println("get thread")
	}

	fmt.Println("start thread")

	s.wait.Add(1)
	go func(f func()) {
		defer s.wait.Done()
		f()
		for {
			select {
			case <-s.ctx.Done():
				<-s.runningThread
				return
			case nextf := <-s.waitingTaskQueue:
				nextf()
				continue
			default:
				<-s.runningThread
				return
			}
		}
	}(f)
	return
}

func (s *waitLimiter) cleanup() {
	s.wait.Wait()
loop:
	for {
		select {
		case s.runningThread <- struct{}{}:
			s.wait.Add(1)
			go func() {
				defer s.wait.Done()
			loop1:
				for {
					select {
					case f := <-s.waitingTaskQueue:
						f()
					default:
						break loop1
					}
				}
			}()
		default:
			break loop
		}
	}
	s.wait.Wait()

}
