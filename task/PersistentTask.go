package task

import (
	"fmt"
	"log"

	"time"
)

type PersistentTaskForm struct {
	ExecutionInterval time.Duration
	Mission           func() error
	TryAgainTimeout   time.Duration
	MaxTryTimes       int
}

func NewPersistentTask(t PersistentTaskForm) *PersistentTask {
	if t.ExecutionInterval == 0 {
		panic("the timeout should not be empty")
	}
	return &PersistentTask{ptf: t}
}

type PersistentTask struct {
	ptf    PersistentTaskForm
	timmer *time.Timer
	//lock     sync.Mutex
	hasNew   bool
	stopChan chan struct{}
}

// 啟動
func (p *PersistentTask) Start() {
	//fmt.Println("PersistentTask Start")
	//fmt.Println("p.ptf.ExecutionInterval: ", p.ptf.ExecutionInterval)
	p.timmer = time.NewTimer(p.ptf.ExecutionInterval)
	//p.lock = sync.Mutex{}
	p.stopChan = make(chan struct{})
	p.hasNew = true
	go func() {
	loop:
		for {
			select {
			case <-p.timmer.C:
				//fmt.Println("PersistentTask time up")
				if !p.hasNew {
					continue
				}
				err := p.ptf.Mission()
				var remainTryTimes = p.ptf.MaxTryTimes
				for err != nil {
					if remainTryTimes == 0 {
						log.Printf("giveup to do mission after %d times try , err: %s , ", p.ptf.MaxTryTimes, err.Error())
						continue loop
					}
					log.Printf("fail to do mission, err: %s , try again after %s", err.Error(), p.ptf.TryAgainTimeout.String())

					time.Sleep(p.ptf.TryAgainTimeout)
					err = p.ptf.Mission()
					remainTryTimes--
				}
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
	if !p.hasNew {
		p.SkipScheduled()
		p.scheduleNext()
	}
}

func (p *PersistentTask) scheduleNext() {
	p.timmer.Reset(p.ptf.ExecutionInterval)
	p.hasNew = true
}

func (p *PersistentTask) SkipScheduled() {
	if !p.timmer.Stop() {
		select {
		case <-p.timmer.C:
		default:
		}
	}
}
