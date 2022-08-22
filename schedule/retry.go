package schedule

import (
	"context"
	"github.com/meowalien/go-meowalien-lib/errs"
	"time"
)

func Retry(ctx context.Context, retryCount int, retryInterval time.Duration, f func(round int) bool) (err error) {
	maxCount := retryCount
	for {
		round := maxCount - retryCount + 1
		if !f(round) {
			return
		}
		if retryCount <= 1 {
			err = errs.New("failed after %d times retry", round)
			return
		}
		select {
		case <-ctx.Done():
			err = errs.New("context done at %dth retry: %w", round, ctx.Err())
			return
		case <-time.After(retryInterval):
			retryCount--
			continue
		}
	}
}
