package task

import (
	"fmt"
	"testing"
	"time"
)

func job() *LinkedTask {
	fmt.Println("ssssssss")
	// do something
	return &LinkedTask{
		Job:       job,
		StartTime: time.Now().Add(time.Microsecond * 100),
	}
}

func TestLinkedTask(t *testing.T) {
	s := NewLinkedTaskScheduler()
	s.ScheduleIfCloser(&LinkedTask{
		Job:       job,
		StartTime: time.Now().Add(time.Microsecond * 100),
	})
}
