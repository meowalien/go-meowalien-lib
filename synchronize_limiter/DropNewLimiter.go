package synchronize_limiter

import (
	"context"
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync"
	"time"
)

type dropNewLimiter struct {
	wait             sync.WaitGroup
	cancel           context.CancelFunc
	ctx              context.Context
	waitingTaskQueue chan func()
	runningThread    chan struct{}
}

func (s *dropNewLimiter) Do(ctx context.Context, f func()) (err error) {
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
	default:
		err = DropMission
		return
	}
	fmt.Println("start thread")
	s.wait.Add(1)
	go func(f func()) {
		f()
		for {
			select {
			case <-s.ctx.Done():
				err = errs.New("limiter stopping")
				fmt.Println(err)
				return
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					err = errs.New("timeout")
				} else {
					err = errs.New("limiter context done: %w", ctx)
				}
				fmt.Println("limiter context done: %w", ctx)
				return
			case nextf := <-s.waitingTaskQueue:
				fmt.Println("run from waiting queue")
				nextf()
				fmt.Println("end run from waiting queue")
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

func (s *dropNewLimiter) Stop() {
	s.cancel()
	s.Wait()
}

func (s *dropNewLimiter) Wait() {
	var c func()
	s.ctx, c = context.WithTimeout(context.Background(), time.Second*10)
	defer c()

	s.wait.Wait()
}
