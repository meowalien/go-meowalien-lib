package arangodb_wrapper

import (
	"context"
	"github.com/arangodb/go-driver"
	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/meowalien/go-meowalien-lib/schedule"

	"time"
)

type RetryWrapper struct {
	driver.Database
	RetryCount    int
	RetryInterval time.Duration
}

func (r *RetryWrapper) Transaction(ctx context.Context, action string, options *driver.TransactionOptions) (res interface{}, err error) {
	errRetry := schedule.Retry(ctx, r.RetryCount, r.RetryInterval, func(_ int) bool {
		var err1 error
		res, err1 = r.Database.Transaction(ctx, action, options)
		if err1 != nil {
			err = errs.New(err, err1)
			return true
		}
		return false
	})
	if errRetry != nil {
		err = errs.New(err, errRetry)
		return
	} else {
		err = nil
	}
	return
}

func (r RetryWrapper) Wrap(db driver.Database) driver.Database {
	return &RetryWrapper{
		Database: db,
	}
}
