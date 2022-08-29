package websocket1

import (
	"context"
	"github.com/meowalien/go-meowalien-lib/errs"
	"time"
)

/*
https://www.rfc-editor.org/rfc/rfc6455
*/

type ConnectionKeeper interface {
	Start(ctx context.Context)
}

type Connection interface {
	OnError(keeper ConnectionKeeper, err error)
	Allocator
	Read(ctx context.Context) (msgType MessageType, data []byte, err error)
}
type Allocator interface {
	Output(ctx context.Context, keeper ConnectionKeeper, msg Message)
	Income(ctx context.Context, keeper ConnectionKeeper, msg Message)
}

const (
	DefaultIncomeBufferSize = 100
	DefaultOutputBufferSize = 100
	DefaultReadTimeout      = time.Second * 5
)

var (
	ErrIncomeQueueFull = errs.New("income queue is full")
)

type Config func(c *config)

func defaultConfig() config {
	return config{
		readTimeout:      DefaultReadTimeout,
		incomeBufferSize: DefaultIncomeBufferSize,
		outputBufferSize: DefaultOutputBufferSize,
	}
}

type config struct {
	readTimeout      time.Duration
	incomeBufferSize int
	outputBufferSize int
}

func NewConnectionKeeper(cnn Connection, configs ...Config) (keeper ConnectionKeeper) {
	conf := defaultConfig()
	for c := range configs {
		configs[c](&conf)
	}
	return &connectionKeeper{
		readTimeout: conf.readTimeout,
		incomeQueue: make(chan Message, conf.incomeBufferSize),
		outputQueue: make(chan Message, conf.outputBufferSize),
		cnn:         cnn,
	}
}

type connectionKeeper struct {
	incomeQueue chan Message
	outputQueue chan Message
	cnn         Connection
	readTimeout time.Duration
}

func (c *connectionKeeper) Start(ctx context.Context) {
	go c.readPump(ctx)
	go c.writePump(ctx)
}

func (c *connectionKeeper) readPump(ctx context.Context) {
	go func() {
		for {
			c.cnn.Income(ctx, c, message)
		}
	}()

	timer := time.NewTimer(0)
	// make sure timer is clean
	<-timer.C
	for {
		select {
		case <-ctx.Done():
		default:
			msgtype, data, err := c.cnn.Read(ctx)
			if err != nil {
				c.cnn.OnError(c, errs.New("error when reading:%w , data:%v , message type: %v", err, data, msgtype))
				continue
			}
			var message Message
			switch msgtype {
			case MessageBinary:
				message = NewBinaryMessage(data)
			case MessageText:
				message = NewTextMessage(data)
			}
			timer.Reset(c.readTimeout)
			select {
			case <-ctx.Done():
			case c.incomeQueue <- message:
			case <-timer.C:
				c.cnn.OnError(c, errs.New("read timeout"))
			}
		}
	}
}

func (c *connectionKeeper) writePump(ctx context.Context) {

}
