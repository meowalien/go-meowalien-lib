package graceful_shutdown

import "fmt"

type GracefulShutdownLevel interface {
	Done() (chFc <-chan func())
	Close()
	fmt.Stringer
	NextLevel(name string) *gracefulShutdownLevel
}

func NewRootLevel(name string) *gracefulShutdownLevel {
	return &gracefulShutdownLevel{contextGroup: newContextGroup(nil), name: name}
}

type gracefulShutdownLevel struct {
	*contextGroup
	name  string
	level uint
}

func (g *gracefulShutdownLevel) String() string {
	return g.name
}

func (g *gracefulShutdownLevel) Level() uint {
	return g.level
}

func (g *gracefulShutdownLevel) NextLevel(name string) *gracefulShutdownLevel {
	return &gracefulShutdownLevel{contextGroup: g.childGroup(), name: name, level: g.level + 1}
}
