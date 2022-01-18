package task

import (
"time"
)

type NewTaskFunc func()(task Task, err error)

type TaskRunner func(interface{}) error

type ErrLogger interface {
	Errorf(template string, args ...interface{})
}

type ContinueTaskForm struct {
	PopUp         NewTaskFunc
	TryAgainAfter time.Duration
	Consumer      TaskRunner
	ErrLogger     ErrLogger
}

//type PopUpTask interface {
//	Start()
//}

func NewPopUpTask(t ContinueTaskForm) *continueTask {
	return &continueTask{t}
}

type Task struct {
	ExecutionTime time.Time `json:"execution_time"`
	Data interface{}
}

type continueTask struct {
	ContinueTaskForm
}

func (c *continueTask) Start() {
	nTask ,err := c.PopUp()
	for err != nil {
		c.ErrLogger.Errorf("error when PopUp tx: %s, try again after %f sec ...", err.Error(), c.TryAgainAfter.Seconds())
		<-time.NewTimer(c.TryAgainAfter).C
		nTask ,err = c.PopUp()
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

	err = c.Consumer(nTask.Data)
	for err != nil {
		c.ErrLogger.Errorf("error when DoTask tx: %s, try again after %f sec ...", err.Error(), c.TryAgainAfter.Seconds())
		<-time.NewTimer(c.TryAgainAfter).C
		err = c.Consumer(nTask.Data)
	}
	c.Start()
}


