package chan_context

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestChanContext(t *testing.T) {
	wg := sync.WaitGroup{}
	level := 200
	childCount := 200
	childCtxCount := 200
	delayRange := 2000
	gp := RootContextGroup("root")

	for i := 0; i < level; i++ {
		if i < childCount {
			gpChild := gp.Child(fmt.Sprintf("child_%d", i))

			ctx, _ := gpChild.NewContext()
			for ii := 0; ii < childCtxCount; ii++ {
				wg.Add(1)
				go func(gp WaitContext, i int, ii int) {
					select {
					case okFc := <-ctx.DonePromise():
						fmt.Printf("done_child_%d_%d\n", i, ii)
						okFc()
					case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
					}
					wg.Done()
				}(ctx, i, ii)
			}
		}

		ctx, _ := gp.NewContext()
		wg.Add(1)

		go func(gp WaitContext, i int) {
			select {
			case okFc := <-ctx.DonePromise():
				fmt.Println("done_root", i)
				okFc()
			case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
			}
			wg.Done()
		}(ctx, i)
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(delayRange)))
	fmt.Println("========================================================")
	gp.Close()
	fmt.Println("after cancel")
	wg.Wait()

}
