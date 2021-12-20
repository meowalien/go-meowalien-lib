package task

import (
"time"
)

type NewTaskFunc func()(task Mission, err error)

type TaskRunner func(interface{}) error

type ErrLogger interface {
	Errorf(template string, args ...interface{})
}

type NewTaskForm struct {
	NewTaskFunc   NewTaskFunc
	TryAgainAfter time.Duration
	DoTask        TaskRunner
	ErrLogger     ErrLogger
}

type ContinueTask interface {
	Start()
}

func NewTaskHolder(t NewTaskForm) ContinueTask {
	return &continueTask{t}
}

type Mission struct {
	ExecutionTime time.Time `json:"execution_time"`
	Data interface{}
}

type continueTask struct {
	NewTaskForm
}

func (c *continueTask) Start() {
	nTask ,err := c.NewTaskFunc()
	for err != nil {
		c.ErrLogger.Errorf("error when NewTaskFunc tx: %s, try again after %f sec ...", err.Error(), c.TryAgainAfter.Seconds())
		<-time.NewTimer(c.TryAgainAfter).C
		nTask ,err = c.NewTaskFunc()
	}
	now := time.Now()

	var scheduleAfter time.Duration

	if nTask.ExecutionTime.Before(now) {
		scheduleAfter = 0
	} else {
		scheduleAfter = nTask.ExecutionTime.Sub(now)
	}

	if scheduleAfter != 0 {
		<-time.NewTimer(scheduleAfter).C
	}

	err = c.DoTask(nTask.Data)
	for err != nil {
		c.ErrLogger.Errorf("error when DoTask tx: %s, try again after %f sec ...", err.Error(), c.TryAgainAfter.Seconds())
		<-time.NewTimer(c.TryAgainAfter).C
		err = c.DoTask(nTask.Data)
	}
	c.Start()
}


