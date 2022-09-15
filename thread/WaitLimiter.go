package thread

import (
	"context"
	"errors"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync"
	"sync/atomic"
)

func newWaitLimiter(cf Config) *waitLimiter {
	ctx, cancel := context.WithCancel(context.Background())
	l := &waitLimiter{
		cond:             sync.Cond{L: &sync.Mutex{}},
		stopChan:         make(chan struct{}),
		cancel:           cancel,
		ctx:              ctx,
		waitingTaskQueue: make(chan func(), cf.WaitingQueueCapacity),
		threadCount:      cf.RunningThreadAmount,
	}
	l.startConsumer()
	return l
}

type waitLimiter struct {
	waitThread          sync.WaitGroup
	waitingTaskQueue    chan func()
	stopChan            chan struct{}
	threadCount         int
	ctx                 context.Context
	cancel              context.CancelFunc
	cond                sync.Cond
	sequentialTaskCount uint64
}

func (s *waitLimiter) WaitQueueClean() {
	s.cond.L.Lock()
	for len(s.waitingTaskQueue) != 0 || s.sequentialTaskCount != 0 {
		s.cond.Wait()
	}
	s.cond.L.Unlock()
}

func (s *waitLimiter) Stop(ctx context.Context) {
	s.cancel()
	s.WaitQueueClean()
	close(s.stopChan)
	s.waitThread.Wait() // wait for all thread closed
}

func (s *waitLimiter) do(ctx context.Context, f func()) (err error) {
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
func (s *waitLimiter) Do(ctx context.Context, ff ...func()) (err error) {
	if len(ff) == 1 {
		return s.do(ctx, ff[0])
	}
	atomic.AddUint64(&s.sequentialTaskCount, uint64(len(ff)))
	f := s.makeFunc(ctx, ff, 0)
	err = s.do(ctx, f)
	return
}

func (s *waitLimiter) makeFunc(ctx context.Context, ff []func(), i int) func() {
	return func() {
		defer atomic.AddUint64(&s.sequentialTaskCount, ^uint64(0))
		ff[i]()
		if i+1 < len(ff) {
			s.waitingTaskQueue <- s.makeFunc(ctx, ff, i+1)
		}
	}
}
func (s *waitLimiter) startConsumer() {
	s.waitThread.Add(s.threadCount)
	for i := 0; i < s.threadCount; i++ {
		go func(i int) {
			defer s.waitThread.Done()
			for {
				select {
				case <-s.stopChan:
					return
				case f := <-s.waitingTaskQueue:
					f()
					s.cond.L.Lock()
					if len(s.waitingTaskQueue) == 0 {
						s.cond.Broadcast()
					}
					s.cond.L.Unlock()
					continue
				}
			}
		}(i)
	}
}
