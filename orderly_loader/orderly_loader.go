package orderly_loader

import (
	"context"
	"github.com/meowalien/go-meowalien-lib/errs"
	"sync"
)

type getter[T any] func(ctx context.Context) (conn T, err error)

type OrderlyLoader[T any] interface {
	SetObject(object *T)
	Loader() (get getter[T], release func())
	IsEmpty() bool
}

func New[T any]() OrderlyLoader[T] {
	return &orderlyLoader[T]{
		cond:              sync.NewCond(&sync.Mutex{}),
		snapshotHoldIndex: 1,
	}
}

type orderlyLoader[T any] struct {
	obj               *T
	cond              *sync.Cond
	snapshotNewIndex  uint64
	snapshotHoldIndex uint64
}

func (s *orderlyLoader[T]) IsEmpty() bool {
	return s.obj == nil
}

func (s *orderlyLoader[T]) SetObject(object *T) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	s.obj = object
}

func (s *orderlyLoader[T]) Loader() (get getter[T], release func()) {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	s.snapshotNewIndex += 1
	index := s.snapshotNewIndex
	var released bool
	onece := sync.Once{}
	release = func() {
		onece.Do(func() {
			s.snapshotHoldIndex += 1
			released = true
			s.cond.L.Unlock()
			s.cond.Broadcast()
		})
	}
	get = func(ctx context.Context) (conn T, err error) {
		s.cond.L.Lock()
		for s.snapshotHoldIndex != index {
			doneChan := make(chan struct{}, 1)
			go func() {
				select {
				case <-ctx.Done():
					err = errs.New("context done in nextOrderlyLoader-getter: %w", ctx.Err())
					release()
				case <-doneChan:
				}
			}()

			s.cond.Wait()
			if released {
				return
			}
			doneChan <- struct{}{}
		}
		conn = *s.obj
		return
	}
	return
}
