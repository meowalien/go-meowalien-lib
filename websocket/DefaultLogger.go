package websocket

import "fmt"

type defaultLogger struct{}

func (d defaultLogger) Debugf(format string, args ...interface{}) {
	format += "\n"
	fmt.Printf(format, args...)
}

func (d defaultLogger) Infof(format string, args ...interface{}) {
	format += "\n"
	fmt.Printf(format, args...)
}

func (d defaultLogger) Warnf(format string, args ...interface{}) {
	format += "\n"
	fmt.Printf(format, args...)
}

func (d defaultLogger) Errorf(format string, args ...interface{}) {
	format += "\n"
	fmt.Printf(format, args...)
}
