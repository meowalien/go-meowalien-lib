package task

import (
	"time"
)

type NewTaskFunc func() (task Task, err error)

type TaskRunner func(interface{}) error

type ErrLogger interface {
	Errorf(template string, args ...interface{})
}

type NewPopUpTaskForm struct {
	Provider      NewTaskFunc
	TryAgainAfter time.Duration
	Consumer      TaskRunner
	ErrLogger             ErrLogger
	GiveUpAfterRetryTimes int
}

//type PopUpTask interface {
//	Start()
//}

func NewPopUpTask(t NewPopUpTaskForm) *continueTask {
	return &continueTask{t}
}

type Task struct {
	ExecutionTime *time.Time `json:"execution_time"`
	Data          interface{}
}

type continueTask struct {
	NewPopUpTaskForm
}

func (c *continueTask) Start() {
	remainRetryTimes := c.GiveUpAfterRetryTimes
	nTask, err := c.Provider()
	for err != nil {
		c.ErrLogger.Errorf("error when Provider tx: %s, try again after %f sec ...", err.Error(), c.TryAgainAfter.Seconds())
		if remainRetryTimes != -1 {
			if remainRetryTimes == 0 {
				c.ErrLogger.Errorf("give up to do Provider after %d times try", c.GiveUpAfterRetryTimes)
				return
			} else {
				remainRetryTimes--
			}
		}

		<-time.NewTimer(c.TryAgainAfter).C
		nTask, err = c.Provider()
	}
	now := time.Now()

	var scheduleAfter time.Duration

	if nTask.ExecutionTime == nil || nTask.ExecutionTime.Before(now) {
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
		if remainRetryTimes != -1 {
			if remainRetryTimes == 0 {
				c.ErrLogger.Errorf("give up to do Provider after %d times try", c.GiveUpAfterRetryTimes)
				return
			} else {
				remainRetryTimes--
			}
		}
		<-time.NewTimer(c.TryAgainAfter).C
		err = c.Consumer(nTask.Data)
	}
	c.Start()
}
