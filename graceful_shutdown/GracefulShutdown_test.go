package graceful_shutdown

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var (
	LevelRoot = NewRootLevel("levelRoot")
	Level1    = LevelRoot.NextLevel("AAA")
	Level2    = Level1.NextLevel("BBB")
)

func TestChanContext(t *testing.T) {
	magnification := 10 * 1
	level := 2 * magnification
	childCount := 2 * magnification
	childCtxCount := 2 * magnification
	delayRange := 200 * magnification

	for i := 0; i < level; i++ {
		if i < childCount {
			for ii := 0; ii < childCtxCount; ii++ {
				go func(i int, ii int) {
					//fmt.Printf("start_%s_%d_%d_%d\n", Level1, Level1.Level(), i, ii)
					select {
					case okFc := <-Level1.PromiseDone():
						fmt.Printf("done_%s_%d_%d_%d\n", Level1, Level1.Level(), i, ii)
						okFc()
					case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
						fmt.Printf("exec_%s_%d_%d_%d\n", Level1, Level1.Level(), i, ii)
					}
				}(i, ii)

				go func(i int, ii int) {
					//fmt.Printf("start_%s_%d_%d_%d\n", Level2, Level2.Level(), i, ii)
					select {
					case okFc := <-Level2.PromiseDone():
						fmt.Printf("done_%s_%d_%d_%d\n", Level2, Level2.Level(), i, ii)
						okFc()
					case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
						fmt.Printf("exec_%s_%d_%d_%d\n", Level2, Level2.Level(), i, ii)
					}
				}(i, ii)

			}
		}

		go func(i int) {
			select {
			case okFc := <-LevelRoot.PromiseDone():
				fmt.Printf("done_%s_%d\n", LevelRoot, i)
				okFc()
			case <-time.After(time.Millisecond * time.Duration(rand.Intn(delayRange))):
				fmt.Printf("exec_%s_%d\n", LevelRoot, i)
			}
		}(i)
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(delayRange)))
	fmt.Println("========================================================")
	LevelRoot.Close()
	fmt.Println("after cancel")

}
