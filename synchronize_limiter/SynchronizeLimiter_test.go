package synchronize_limiter

import (
	"context"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//defer cancel()
	l := NewLimiter(Config{
		QueueFullStrategy:  Strategy_Wait,
		WaitingQueueLimit:  10,
		RunningThreadLimit: 1,
	})
	for i := 0; i < 100; i++ {
		ii := i
		fmt.Println("put : ", ii)
		ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
		err := l.Do(ctx, func() {
			fmt.Println("ST: ", ii)
			//if ii%2 == 0 {
			//	time.Sleep(time.Second / 4)
			//} else {
			time.Sleep(time.Second / 2)
			//}
			fmt.Println("EDT: ", ii)
		})
		if err != nil {
			err = errs.New(err)
			//panic(err)
			fmt.Println(err)
		}
	}
	time.Sleep(time.Second * 2)
	fmt.Println("\t\tstart to wait")
	l.Stop()

}
