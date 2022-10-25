package schedule

import (
	"context"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"testing"
	"time"
)

func TestRetryFail(t *testing.T) {
	err := Retry(context.TODO(), 5, time.Millisecond*300, func(round int) bool {
		fmt.Println("aaaaaaa: ", round)
		return true
	})
	fmt.Println(err)
}

func TestRetryDone(t *testing.T) {
	err := Retry(context.TODO(), 5, time.Millisecond*300, func(round int) bool {
		fmt.Println("aaaaaaa: ", round)
		return round != 3
	})
	fmt.Println(err)
}

func TestRetryErr(t *testing.T) {
	var err error
	errRetry := Retry(context.TODO(), 5, time.Millisecond*300, func(round int) bool {
		err = errs.New(err, fmt.Errorf("error %d", round))
		return true
	})
	if !errRetry {
		err = errs.New("not done")

	} else {
		err = nil
	}

	fmt.Println(err)
}
