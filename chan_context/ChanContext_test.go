package chan_context

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type GracefulShutdownLevel struct {
	GroupContext[uint8]
	name string
}

func (g *GracefulShutdownLevel) String() string {
	return g.name
}

func ChildLevel(g *GracefulShutdownLevel, name string) *GracefulShutdownLevel {
	return &GracefulShutdownLevel{GroupContext: g.Child(g.Key() + 1), name: name}
}

var (
	LevelRoot = &GracefulShutdownLevel{GroupContext: NewContextGroup(uint8(0)), name: "levelRoot"}
	Level1    = ChildLevel(LevelRoot, "level1")
	Level2    = ChildLevel(Level1, "level2")
)

func TestChanContext(t *testing.T) {
	wg := sync.WaitGroup{}
	level := 200
	childCount := 200
	childCtxCount := 200
	delayRange := 20000

	for i := 0; i < level; i++ {
		if i < childCount {
			ctx, _ := Level1.Context()
			ctx1, _ := Level2.Context()
			for ii := 0; ii < childCtxCount; ii++ {
				wg.Add(1)
				go func(gp WaitContext, i int, ii int) {
					select {
					case okFc := <-ctx.DonePromise():
						fmt.Printf("done_%s_%d\n", Level1, ii)
						okFc()
					case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
						fmt.Printf("exec_%s_%d\n", Level1, ii)
					}
					wg.Done()
				}(ctx, i, ii)

				wg.Add(1)
				go func(gp WaitContext, i int, ii int) {
					select {
					case okFc := <-ctx1.DonePromise():
						fmt.Printf("done_%s_%d\n", Level2, ii)
						okFc()
					case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
						fmt.Printf("exec_%s_%d\n", Level2, ii)
					}
					wg.Done()
				}(ctx1, i, ii)
			}
		}

		ctx, _ := LevelRoot.Context()
		wg.Add(1)

		go func(gp WaitContext, i int) {
			select {
			case okFc := <-ctx.DonePromise():
				fmt.Printf("done_%s_%d\n", LevelRoot, i)
				okFc()
			case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
				fmt.Printf("exec_%s_%d\n", LevelRoot, i)
			}
			wg.Done()
		}(ctx, i)
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(delayRange)))
	fmt.Println("========================================================")
	LevelRoot.Close()
	fmt.Println("after cancel")
	wg.Wait()

}
