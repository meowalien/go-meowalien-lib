package thread

import (
	"context"
	"errors"
	"github.com/meowalien/go-meowalien-lib/errs"
)

type Stop interface {
	Stop(ctx context.Context)
}

type SynchronizeLimiter interface {
	Do(ctx context.Context, f func()) (err error)
	Stop
}

type Strategy int

const (
	Strategy_Wait = iota // default
)

type Config struct {
	QueueFullStrategy  Strategy
	WaitingQueueLimit  int
	RunningThreadLimit int
	Ctx                context.Context
}

var DropMission = errors.New("drop")

func NewSynchronizeLimiter(cf Config) SynchronizeLimiter {
	if cf.Ctx == nil {
		cf.Ctx = context.TODO()
	}
	if cf.RunningThreadLimit < 1 {
		panic("running thread limit is less than 1")
	}

	switch cf.QueueFullStrategy {
	case Strategy_Wait:
		ctx, cancel := context.WithCancel(cf.Ctx)
		l := &waitLimiter{
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
