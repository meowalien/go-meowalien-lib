package schedule

import (
	"fmt"
	"time"
)

type Timer struct {
	*time.Timer
	Canceled chan time.Time
}

func (rf *Timer) Cancel() {
	fmt.Println("Timer-Cancel")
	if !rf.Stop() {
		select {
		case <-rf.Timer.C:
		default:
		}
	}

	rf.Canceled <- time.Now()
	fmt.Println("Timer-Cancel-send")
}

func NewCancelableTimer(duration time.Duration) *Timer {
	return &Timer{
		Timer:    time.NewTimer(duration),
		Canceled: make(chan time.Time, 1),
	}
}
