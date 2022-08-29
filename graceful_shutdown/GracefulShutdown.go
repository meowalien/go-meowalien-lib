package graceful_shutdown

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/contexts"
)

type GracefulShutdownLevel interface {
	PromiseDone() (chFc <-chan func())
	Close()
	fmt.Stringer
	NextLevel(name string) *gracefulShutdownLevel
}

func NewRootLevel(name string) *gracefulShutdownLevel {
	return &gracefulShutdownLevel{ContextGroup: contexts.NewContextGroup(nil), name: name}
}

type gracefulShutdownLevel struct {
	contexts.ContextGroup
	name  string
	level uint
}

func (g *gracefulShutdownLevel) String() string {
	if g.name == "" {
		return fmt.Sprintf("Level%d", g.level)
	}
	return g.name
}

func (g *gracefulShutdownLevel) Level() uint {
	return g.level
}

func (g *gracefulShutdownLevel) NextLevel(name string) *gracefulShutdownLevel {
	return &gracefulShutdownLevel{ContextGroup: g.ChildGroup(), name: name, level: g.level + 1}
}
