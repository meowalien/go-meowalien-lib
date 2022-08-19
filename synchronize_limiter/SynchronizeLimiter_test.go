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
	limiter := NewLimiter(Config{
		QueueFullStrategy:  Strategy_Wait,
		WaitingQueueLimit:  1024,
		RunningThreadLimit: 1,
	})
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	go func() {
		for i := 0; i < 10000; i++ {
			ii := i
			ctx, _ := context.WithTimeout(context.Background(), time.Second*100)
			//wg.Add(1)
			err := limiter.Do(ctx, func() {
				//wg.Done()
				fmt.Println("AAAAAAAAAAAAAAAAAAA: ", ii)
				time.Sleep(time.Microsecond * 20)
				//l.snapshot()
				//cancel()
			})
			if err != nil {
				err = errs.New(err)
				fmt.Println(err)
				return
			}
		}
	}()

	//}()
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	for i := 0; i < 10000; i++ {
	//		ii := i
	//		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	//		wg.Add(1)
	//		err := limiter.Do(ctx, func() {
	//			wg.Done()
	//			fmt.Println("AAAAAAAAAAAAAAAAAAA: ", ii)
	//			//time.Sleep(time.Microsecond * 26)
	//			//l.snapshot()
	//			//cancel()
	//		})
	//		if err != nil {
	//			err = errs.New(err)
	//			fmt.Println(err)
	//		}
	//	}
	//}()
	time.Sleep(time.Second)
	fmt.Println("\t\tstart to wait")
	limiter.Stop(context.TODO())
	//wg.Wait()

}
