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
	p.timmer = time.NewTimer(p.ptf.TimeOut)
	p.lock = sync.Mutex{}
	p.stopChan = make(chan struct{}, 0)

loop:
	for {
		select {
		case <-p.timmer.C:
			if !p.hasNew {
				continue
			}
			p.ptf.Mission()
			p.hasNew = false
			p.timmer.Reset(p.ptf.TimeOut)
		case <-p.stopChan:
			if !p.timmer.Stop() {
				<-p.timmer.C
			}
			break loop
		}
	}
	fmt.Println("persistentTask stop")
}

func (p *PersistentTask) Stop() {
	close(p.stopChan)
}

// 通知有新任務
func (p *PersistentTask) Active() {
	if !p.hasNew {
		p.hasNew = true
	}
}
