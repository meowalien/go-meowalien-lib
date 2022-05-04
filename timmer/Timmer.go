package timmer

import "time"
type Timer struct {
	*time.Timer
	Canceled chan time.Time
}

func (rf *Timer) Cancel() {
	if rf.Timer != nil {
		if !rf.Timer.Stop() {
			select {
			case <-rf.Timer.C:
			default:
			}
		}
		rf.Canceled <- time.Now()
	}
}

func NewCancelableTimer(duration time.Duration) *Timer {
	return &Timer{
		Timer: time.NewTimer(duration),
		Canceled: make(chan time.Time, 1),
	}
}

