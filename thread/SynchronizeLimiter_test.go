package thread

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
	limiter := NewSynchronizeLimiter(Config{
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

func TestLiniter(t *testing.T) {
	limiter := NewSynchronizeLimiter(Config{
		QueueFullStrategy:  Strategy_Wait,
		WaitingQueueLimit:  1024,
		RunningThreadLimit: 10,
	})

	err := limiter.Do(context.TODO(), func() {
		fmt.Println("AAAAAAAAAAAAAAAAAAA: ")
		time.Sleep(time.Second * 2)
		fmt.Println("AAAAAAAAAAAAAAAAAAA-afterwait: ")
	}, func() {
		fmt.Println("BBBBBBBBBBBBBBBBBBBB: ")
		time.Sleep(time.Second)
		fmt.Println("BBBBBBBBBBBBBBBBBBBB-afterwait: ")

	}, func() {
		fmt.Println("CCCCCCCCCCCCCCCCCCCC: ")
		time.Sleep(time.Second / 2)
		fmt.Println("CCCCCCCCCCCCCCCCCCCC-afterwait: ")

	})
	if err != nil {
		panic(err)
	}
	err = limiter.Do(context.TODO(), func() {
		fmt.Println("ADAADDADADADADADADADAD: ")
		time.Sleep(time.Second * 2)
		fmt.Println("ADAADDADADADADADADADAD-afterwait: ")
	}, func() {
		fmt.Println("BVBVBVBVBVBVBVBVBVBVBVB: ")
		time.Sleep(time.Second)
		fmt.Println("BVBVBVBVBVBVBVBVBVBVBVB-afterwait: ")

	}, func() {
		fmt.Println("CXCXCXCXCXCXCXCXCXCXc: ")
		time.Sleep(time.Second / 2)
		fmt.Println("CXCXCXCXCXCXCXCXCXCXc-afterwait: ")

	})
	if err != nil {
		panic(err)
	}
	//time.Sleep(time.Second * 1)
	//limiter.Stop(context.TODO())
	limiter.WaitQueueClean()
	fmt.Println(err)
}
