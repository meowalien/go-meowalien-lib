package synchronize_limiter

import (
	"context"
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync"
)

type Stop interface {
	Stop()
}

type Wait interface {
	Wait()
}

type Limiter interface {
	Do(ctx context.Context, f func(ctx context.Context)) (err error)
	Stop
	Wait
}

type Config struct {
	RunningThreadLimit int
	Ctx                context.Context
	WaitingQueueLimit  int
}

func NewLimiter(cf Config) Limiter {
	if cf.Ctx == nil {
		cf.Ctx = context.TODO()
	}
	if cf.RunningThreadLimit < 1 {
		panic("running thread limit is less than 1")
	}
	ctx, cancel := context.WithCancel(cf.Ctx)
	return &limiter{
		cancel:        cancel,
		ctx:           ctx,
		waitingTask:   make(chan func(), cf.WaitingQueueLimit),
		runningThread: make(chan struct{}, cf.RunningThreadLimit),
	}
}

type limiter struct {
	waitingTask   chan func()
	runningThread chan struct{}
	wait          sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
}

func (s *limiter) Stop() {
	s.cancel()
	s.Wait()
}

func (s *limiter) Wait() {
	s.wait.Wait()
}

func (s *limiter) Do(ctx context.Context, f func(ctx context.Context)) (err error) {
	select {
	case s.waitingTask <- func() {
		f(ctx)
	}:
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
	go func(f func(ctx context.Context)) {
		f(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case nextf := <-s.waitingTask:
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
