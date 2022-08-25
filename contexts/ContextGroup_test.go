package contexts

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type GracefulShutdownLevel struct {
	ContextGroup[uint8]
	name string
}

func (g *GracefulShutdownLevel) String() string {
	return g.name
}

func ChildLevel(g *GracefulShutdownLevel, name string) *GracefulShutdownLevel {
	return &GracefulShutdownLevel{ContextGroup: g.Child(g.Key() + 1), name: name}
}

var (
	LevelRoot = &GracefulShutdownLevel{ContextGroup: NewContextGroup(uint8(0)), name: "levelRoot"}
	Level1    = ChildLevel(LevelRoot, "level1")
	Level2    = ChildLevel(Level1, "level2")
)

func TestChanContext(t *testing.T) {
	wg := sync.WaitGroup{}
	level := 20
	childCount := 20
	childCtxCount := 20
	delayRange := 2000

	for i := 0; i < level; i++ {
		if i < childCount {
			for ii := 0; ii < childCtxCount; ii++ {
				wg.Add(1)
				go func(i int, ii int) {
					select {
					case okFc := <-Level1.Done():
						fmt.Printf("done_%s_%d\n", Level1, ii)
						okFc()
					case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
						fmt.Printf("exec_%s_%d\n", Level1, ii)
					}
					wg.Done()
				}(i, ii)

				wg.Add(1)
				go func(i int, ii int) {
					select {
					case okFc := <-Level2.Done():
						fmt.Printf("done_%s_%d\n", Level2, ii)
						okFc()
					case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
						fmt.Printf("exec_%s_%d\n", Level2, ii)
					}
					wg.Done()
				}(i, ii)
			}
		}

		wg.Add(1)

		go func(i int) {
			select {
			case okFc := <-LevelRoot.Done():
				fmt.Printf("done_%s_%d\n", LevelRoot, i)
				okFc()
			case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
				fmt.Printf("exec_%s_%d\n", LevelRoot, i)
			}
			wg.Done()
		}(i)
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(delayRange)))
	fmt.Println("========================================================")
	LevelRoot.Close()
	//cancel()
	fmt.Println("after cancel")
	wg.Wait()

}
