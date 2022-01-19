package task

import (
	"fmt"
	"sync"
	"time"
)

type PersistentTaskForm struct {
	TimeOut time.Duration
	Mission func()
}

func NewPersistentTask(t PersistentTaskForm) *PersistentTask {
	if t.TimeOut == 0 {
		panic("the timeout should not be empty")
	}
	return &PersistentTask{ptf: t}
}

type PersistentTask struct {
	ptf      PersistentTaskForm
	timmer   *time.Timer
	lock     sync.Mutex
	hasNew   bool
	stopChan chan struct{}
}

// 啟動
func (p *PersistentTask) Start() {
	fmt.Println("PersistentTask Start")
	fmt.Println("p.ptf.TimeOut: ", p.ptf.TimeOut)
	p.timmer = time.NewTimer(p.ptf.TimeOut)
	p.lock = sync.Mutex{}
	p.stopChan = make(chan struct{}, 0)
	p.hasNew = true
	go func() {
	loop:
		for {
			select {
			case <-p.timmer.C:
				fmt.Println("time up")
				if !p.hasNew {
					continue
				}
				p.ptf.Mission()
				p.hasNew = false
			case <-p.stopChan:
				if !p.timmer.Stop() {
					<-p.timmer.C
				}
				break loop
			}
		}
		fmt.Println("persistentTask stop")
	}()

}

func (p *PersistentTask) Stop() {
	fmt.Println("enter Stop")

	close(p.stopChan)
}

// 通知有新任務
func (p *PersistentTask) Active() {
	fmt.Println("enter Active")
	defer fmt.Println("end Active")
	if !p.hasNew {
		// if !t.Stop() {
		//  <-t.C
		// }
		// t.Reset(d)
		if !p.timmer.Stop(){
			select {
			case <-p.timmer.C:
			default:
			}
		}
		//fmt.Println("timmerStop: ",timmerStop)

		//if !{
		//

		//}
		resetbool := p.timmer.Reset(p.ptf.TimeOut)
		fmt.Println("Resetbool: ", resetbool)
		p.hasNew = true
	}

}
