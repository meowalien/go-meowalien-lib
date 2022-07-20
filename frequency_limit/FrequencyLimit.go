package frequency_limit

import (
	"sync"
	"time"
)

func New(gap time.Duration) FrequencyLimiter {
	return &frequencyLimit{gap: gap}
}

type FrequencyLimiter interface {
	Do(key interface{}, f func())
	Trigger(key interface{}) bool
}

type frequencyLimit struct {
	m   sync.Map
	gap time.Duration
}

func (f *frequencyLimit) Trigger(key interface{}) bool {
	_, loaded := f.m.LoadOrStore(key, nil)
	time.AfterFunc(f.gap, func() {
		f.m.Delete(key)
	})
	return !loaded
}

func (f *frequencyLimit) Do(key interface{}, fc func()) {
	if f.Trigger(key) {
		fc()
	}
}
