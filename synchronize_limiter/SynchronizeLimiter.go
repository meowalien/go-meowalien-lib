package synchronize_limiter

import (
	"context"
	"errors"
	"github.com/meowalien/go-meowalien-lib/errs"
)

type Stop interface {
	Stop()
}

type Wait interface {
	Wait()
}

type Limiter interface {
	Do(ctx context.Context, f func()) (err error)
	Stop
	Wait
}

type Strategy int

const (
	Strategy_Wait = iota // default
	Strategy_DropNew

	// when Strategy_DropOld, the WaitingQueueLimit should be greater than 0
	Strategy_DropOld
)

type Config struct {
	QueueFullStrategy  Strategy
	WaitingQueueLimit  int
	RunningThreadLimit int
	Ctx                context.Context
}

var DropMission = errors.New("drop")

func NewLimiter(cf Config) Limiter {
	if cf.Ctx == nil {
		cf.Ctx = context.TODO()
	}
	if cf.RunningThreadLimit < 1 {
		panic("running thread limit is less than 1")
	}

	switch cf.QueueFullStrategy {
	case Strategy_Wait:
		ctx, cancel := context.WithCancel(cf.Ctx)
		return &waitLimiter{
			cancel:           cancel,
			ctx:              ctx,
			waitingTaskQueue: make(chan func(), cf.WaitingQueueLimit),
			runningThread:    make(chan struct{}, cf.RunningThreadLimit),
		}
	case Strategy_DropNew:
		ctx, cancel := context.WithCancel(cf.Ctx)
		return &dropNewLimiter{
			cancel:           cancel,
			ctx:              ctx,
			waitingTaskQueue: make(chan func(), cf.WaitingQueueLimit),
			runningThread:    make(chan struct{}, cf.RunningThreadLimit),
		}
	case Strategy_DropOld:
		if cf.WaitingQueueLimit == 0 {
			panic(errs.New("the WaitingQueueLimit should be greater than 0 in Strategy_DropOld"))
		}
		ctx, cancel := context.WithCancel(cf.Ctx)
		return &dropOldLimiter{
			cancel:           cancel,
			ctx:              ctx,
			waitingTaskQueue: make(chan func(), cf.WaitingQueueLimit),
			runningThread:    make(chan struct{}, cf.RunningThreadLimit),
		}
	default:
		panic(errs.New("unsupported queue full strategy: %v", cf.QueueFullStrategy))
	}

}
