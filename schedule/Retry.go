package schedule

import (
	"context"
	"time"
)

func Retry(ctx context.Context, retryCount int, retryInterval time.Duration, f func(round int) bool) (done bool) {
	maxCount := retryCount
	for {
		round := maxCount - retryCount + 1
		if !f(round) {
			done = true
			return
		}
		if retryCount <= 1 {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(retryInterval):
			retryCount--
			continue
		}
	}
}
