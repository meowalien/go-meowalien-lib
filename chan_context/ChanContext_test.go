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
	//wg1 := sync.WaitGroup{}
	level := 10
	childCount := 5
	childCtxCount := 5
	delayRange := 20
	gp := RootContextGroup("root")
	//lock := sync.Mutex{}

	for i := 0; i < level; i++ {
		if i < childCount {
			gpChild := gp.Child(fmt.Sprintf("child_%d", i))

			ctx, _ := gpChild.NewContext()
			for ii := 0; ii < childCtxCount; ii++ {
				wg.Add(1)
				//wg1.Add(1)
				go func(gp WaitContext, i int, ii int) {
					//wg1.Done()
					//fmt.Println("before sleep: ", i)
					doneChan, okFc := ctx.DonePromise()
					defer okFc()
					//time.Sleep()
					fmt.Printf("start_select: %d_%d\n", i, ii)
					select {
					case <-doneChan:
						fmt.Printf("done_child_%d_%d\n", i, ii)
						//time.Sleep(time.Second / 2)

					case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
						fmt.Printf("exec_child_%d_%d\n", i, ii)
					}

					wg.Done()
				}(ctx, i, ii)
			}
		}

		ctx, _ := gp.NewContext()
		wg.Add(1)
		//wg1.Add(1)

		go func(gp WaitContext, i int) {
			//wg1.Done()
			//fmt.Println("before sleep: ", i)
			doneChan, okFc := ctx.DonePromise()
			defer okFc()

			//time.Sleep()
			fmt.Printf("start_select: %d\n", i)
			select {
			case <-doneChan:
				fmt.Println("done_root", i)
				//time.Sleep(time.Second / 2)
				//okFc()
			case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
				fmt.Println("exec_root", i)
			}
			//okFc()
			wg.Done()
		}(ctx, i)
	}

	//time.Sleep(time.Second * 2)
	//wg1.Wait()
	//wg.Add(childCount*childCtxCount + level)
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(delayRange)))
	fmt.Println("========================================================")
	gp.Close()
	//cancel()
	fmt.Println("after cancel")
	wg.Wait()

}
