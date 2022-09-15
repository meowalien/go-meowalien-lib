package thread

import (
	"context"
	"github.com/meowalien/go-meowalien-lib/errs"
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
	if cf.WaitingQueueLimit < 1 {
		panic("waiting queue limit is less than 1")
	}
	switch cf.QueueFullStrategy {
	case Strategy_Wait:
		return newWaitLimiter(cf)
	case Strategy_DropOldest:
		return newDropOldLimiter(cf)

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
	Strategy_Wait       Strategy = 1 // default
	Strategy_DropOldest Strategy = 2
)

type Config struct {
	QueueFullStrategy  Strategy
	WaitingQueueLimit  int
	RunningThreadLimit int
}

//var DropMission = errors.New("drop")
