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
		WaitingQueueLimit:  3,
		RunningThreadLimit: 50,
		Ctx:                context.TODO(),
	})
	for i := 0; i < 100; i++ {
		ii := i
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := l.Do(ctx, func() {
			fmt.Println("ST: ", ii)
			//if ii%2 == 0 {
			//	time.Sleep(time.Second / 4)
			//} else {
			//	time.Sleep(time.Second / 2)
			//}
			fmt.Println("EDT: ", ii)
		})
		cancel()
		if err != nil {
			err = errs.New(err)
			panic(err)
		}
	}
	fmt.Println("start to wait")
	l.Wait()

}
