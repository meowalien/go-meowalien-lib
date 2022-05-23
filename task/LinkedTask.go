package task

import (
	"sync"
	"time"
)

type LinkedTaskScheduler interface {
	ScheduleIfCloser(task *LinkedTask)
}
type LinkedTask struct {
	Job       func() *LinkedTask
	StartTime time.Time
}

func NewLinkedTaskScheduler() LinkedTaskScheduler {
	t := time.NewTimer(0)
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	return &linkedTaskScheduler{timer: t, cancelChan: make(chan struct{}, 1)}
}

type linkedTaskScheduler struct {
	timer       *time.Timer
	cancelChan  chan struct{}
	closestTask *LinkedTask
	scheduling  bool
	lock        sync.Mutex
}

func (l *linkedTaskScheduler) ScheduleIfCloser(task *LinkedTask) {
	// 第一個任務
	if l.closestTask == nil {
		go l.do(task)
		return
	}
	// 新任務更靠近
	if l.closestTask.StartTime.After(task.StartTime) {
		// stop old task and schedule new task
		if l.scheduling {
			l.cancelChan <- struct{}{}
		}
		go l.do(task)
	}

	return
}

// 執行最新任務，如果時間為未來，等待
func (l *linkedTaskScheduler) do(task *LinkedTask) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.closestTask = task
	nowTime := time.Now()
	if nowTime.Before(task.StartTime) {
		l.timer.Reset(task.StartTime.Sub(nowTime))
		l.scheduling = true
		select {
		case <-l.cancelChan:
			if !l.timer.Stop() {
				select {
				case <-l.timer.C:
				default:
				}
			}
			l.scheduling = false
			return
		case <-l.timer.C:
			l.scheduling = false
		}
	}
	newTask := task.Job()
	if newTask == nil {
		return
	}
	l.do(newTask)
	return
}
