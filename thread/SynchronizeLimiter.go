package thread

import (
	"context"
	"errors"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync"
)

type SynchronizeLimiter interface {
	Do(ctx context.Context, f ...func()) (err error)
	Stop
	Wait
}

func NewSynchronizeLimiter(cf Config) SynchronizeLimiter {
	if cf.RunningThreadLimit < 1 {
		panic("running thread limit is less than 1")
	}

	ctx, cancel := context.WithCancel(context.Background())
	switch cf.QueueFullStrategy {
	case Strategy_Wait:
		l := &waitLimiter{
			cond:             sync.Cond{L: &sync.Mutex{}},
			stopChan:         make(chan struct{}),
			cancel:           cancel,
			ctx:              ctx,
			waitingTaskQueue: make(chan func(), cf.WaitingQueueLimit),
			threadCount:      cf.RunningThreadLimit,
		}
		l.startConsumer()
		return l
	default:
		panic(errs.New("unsupported queue full strategy: %v", cf.QueueFullStrategy))
	}
}

type Stop interface {
	Stop(ctx context.Context)
}

type Wait interface {
	WaitQueueClean()
}

type Strategy int

const (
	Strategy_Wait = iota // default
)

type Config struct {
	QueueFullStrategy  Strategy
	WaitingQueueLimit  int
	RunningThreadLimit int
}

var DropMission = errors.New("drop")
