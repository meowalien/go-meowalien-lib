package thread

import (
	"context"
	"errors"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync"
	"sync/atomic"
)

func newDropOldLimiter(cf Config) *dropOldestLimiter {
	ctx, cancel := context.WithCancel(context.Background())

	d := &dropOldestLimiter{
		cond:             sync.Cond{L: &sync.Mutex{}},
		stopChan:         make(chan struct{}),
		cancel:           cancel,
		ctx:              ctx,
		waitingTaskQueue: make(chan func(), cf.WaitingQueueCapacity),
		threadCount:      cf.RunningThreadAmount,
	}
	d.startConsumer()
	return d
}

type dropOldestLimiter struct {
	waitThread          sync.WaitGroup
	threadCount         int
	stopChan            chan struct{}
	waitingTaskQueue    chan func()
	cond                sync.Cond
	cancel              context.CancelFunc
	ctx                 context.Context
	sequentialTaskCount uint64
	doLock              sync.Mutex
}

func (s *dropOldestLimiter) WaitQueueClean() {
	s.cond.L.Lock()
	for len(s.waitingTaskQueue) != 0 || s.sequentialTaskCount != 0 {
		s.cond.Wait()
	}
	s.cond.L.Unlock()
}

func (s *dropOldestLimiter) Stop(ctx context.Context) {
	s.cancel()
	s.WaitQueueClean()
	close(s.stopChan)
	s.waitThread.Wait() // wait for all thread closed
}
func (s *dropOldestLimiter) Do(ctx context.Context, ff ...func()) (err error) {
	// to prevent multiple goroutine call Do() at the same time,
	// so that the s.waitingTaskQueue will only be written by one goroutine
	s.doLock.Lock()
	defer s.doLock.Unlock()
	if len(ff) == 1 {
		return s.do(ctx, ff[0])
	}
	atomic.AddUint64(&s.sequentialTaskCount, uint64(len(ff)))
	f := s.makeFunc(ctx, ff, 0)
	err = s.do(ctx, f)
	return
}

func (s *dropOldestLimiter) do(ctx context.Context, f func()) (err error) {
begin:
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
	default:
		// to prevent the other goroutine from testing len(s.waitingTaskQueue)==0
		// before the f is written into s.waitingTaskQueue
		s.cond.L.Lock()
		defer s.cond.L.Unlock()

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
		case <-s.waitingTaskQueue:
			//	consume oldest
		default:
			//	nothing to consume
		}

		goto begin
	}
}

func (s *dropOldestLimiter) makeFunc(ctx context.Context, ff []func(), i int) func() {
	return func() {
		defer atomic.AddUint64(&s.sequentialTaskCount, ^uint64(0))
		ff[i]()
		if i+1 < len(ff) {
			s.waitingTaskQueue <- s.makeFunc(ctx, ff, i+1)
		}
	}
}
func (s *dropOldestLimiter) startConsumer() {
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
