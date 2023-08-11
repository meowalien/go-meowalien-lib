package graceful_shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var afterAllStop = make(chan struct{})

var c = make(chan os.Signal, 1)

func init() {
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		signal.Stop(c)
		runStopStack()
	}()
}

func StopNow(ctx context.Context) error {
	select {
	case <-ctx.Done():
		// abort
		return ctx.Err()
	case <-afterAllStop:
		// already stopped
		return nil
	case c <- syscall.SIGTERM:
		// successfully sent
		select {
		case <-ctx.Done():
			// abort
			return ctx.Err()
		case <-afterAllStop:
			// successfully stopped
			return nil
		}
	default:
		// already notified, wait for stop
		select {
		case <-ctx.Done():
			// abort
			return ctx.Err()
		case <-afterAllStop:
			// successfully stopped
			return nil
		}
	}
}

type stopMission struct {
	f    func()
	name string
}

var stopStack = make([]stopMission, 0)

func Add(name string, f func()) {
	stopStack = append(stopStack, stopMission{
		f:    f,
		name: name,
	})
}

func runStopStack() {
	for i := len(stopStack) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic when running stopStack:%s , %+v\n", stopStack[i].name, r)
				}
			}()
			stopStack[i].f()
		}()
	}
}

func WaitAllStop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-afterAllStop:
		return nil
	}
}
