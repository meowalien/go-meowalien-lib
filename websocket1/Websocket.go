package websocket1

import (
	"context"
	"errors"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/contexts"
	"github.com/meowalien/go-meowalien-lib/errs"
	"io"
	"sync"
	"time"
)

/*
https://www.rfc-editor.org/rfc/rfc6455
*/

type ConnectionKeeper interface {
	MessageSender
	Start()
	Close() error
}
type MessageSender interface {
	SendMessage(ctx context.Context, message Message) (err error)
}
type Allocator func(ctx context.Context, msg Message)
type OnErrorCallback func(keeper ConnectionKeeper, err error)
type WebsocketReader func(ctx context.Context) (msgType MessageType, data []byte, err error)
type WebsocketWriter func(ctx context.Context, typ MessageType, p []byte) (err error)

type ConnectionAdapter struct {
	OnError         OnErrorCallback
	Dispatch        Allocator
	WebsocketReader WebsocketReader
	WebsocketWriter WebsocketWriter
	Close           func() error
}

const (
	DefaultIncomeBufferSize = 100
	DefaultOutputBufferSize = 100
	DefaultReadTimeout      = time.Second * 5
)

var (
	ErrIncomeQueueFull = errs.New("income queue is full")
	ErrOutputQueueFull = errs.New("output queue is full")
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
	pingInterval     time.Duration
	readTimeout      time.Duration
	incomeBufferSize int
	outputBufferSize int
}

func NewConnectionKeeper(cnn ConnectionAdapter, configs ...Config) (keeper ConnectionKeeper) {
	conf := defaultConfig()
	for c := range configs {
		configs[c](&conf)
	}
	onece := sync.Once{}
	cnn.Close = func() (err error) {
		onece.Do(func() {
			err = cnn.Close()
		})
		return
	}
	return &connectionKeeper{
		pingInterval:     conf.pingInterval,
		rootContextGroup: contexts.NewContextGroup(nil),
		readTimeout:      conf.readTimeout,
		incomeQueue:      make(chan Message, conf.incomeBufferSize),
		outputQueue:      make(chan Message, conf.outputBufferSize),
		cnn:              cnn,
	}
}

type connectionKeeper struct {
	incomeQueue      chan Message
	outputQueue      chan Message
	cnn              ConnectionAdapter
	readTimeout      time.Duration
	pingInterval     time.Duration
	rootContextGroup contexts.ContextGroup
}

func (c *connectionKeeper) Close() error {
	// close the readPump, dispatchPump, writePump
	c.rootContextGroup.Close()
	// then close the connection
	return c.cnn.Close()
}

// SendMessage will try to put the message into the output queue,
// if the queue is full, it will wait till the given ctx done and return ErrOutputQueueFull
func (c *connectionKeeper) SendMessage(ctx context.Context, message Message) (err error) {
	select {
	case <-ctx.Done():
		return ErrOutputQueueFull
	case c.outputQueue <- message:
	}
	return
}

func (c *connectionKeeper) Start() {
	// to make sure the readPumpCtx will be done before the writePumpCtx
	writePumpCtx := contexts.NewContextGroup(c.rootContextGroup)
	dispatchPumpCtx := writePumpCtx.ChildGroup()
	readPumpCtx := dispatchPumpCtx.ChildGroup()

	go c.writePump(writePumpCtx)       // close 3st
	go c.dispatchPump(dispatchPumpCtx) // close 2st
	go c.readPump(readPumpCtx)         // close 1st
	<-c.rootContextGroup.Done()
}

// readPump will read the message from the connection and put it into the income queue
func (c *connectionKeeper) readPump(ctx contexts.PromiseContext) {
	timer := time.NewTimer(0)
	// make sure timer is clean before use
	<-timer.C

loop:
	for {
		select {
		case <-ctx.Done():
			// break when the context is done, but the dispatchPump is still running, it will drain the queue
			return
		default:
			msgtype, data, err := c.cnn.WebsocketReader(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					err = c.Close()
					if err != nil {
						err = errs.New(err)
						c.cnn.OnError(c, err)
						return
					}
					return
				}
				c.cnn.OnError(c, errs.New("error when reading:%w , data:%v , message type: %v", err, data, msgtype))
				return
			}
			message := NewMessage(c, msgtype, data)
			fmt.Println("new message:", message)
			timer.Reset(c.readTimeout)
			select {
			case ok := <-ctx.PromiseDone():
				// cancel to push the message to the queue
				if !timer.Stop() {
					<-timer.C
				}
				ok()
				break loop
			case c.incomeQueue <- message:
				if !timer.Stop() {
					<-timer.C
				}
				continue
			case <-timer.C:
				c.cnn.OnError(c, errs.New("%w, dropping message: %v", ErrIncomeQueueFull, message))
				continue
			}
		}
	}
}

func (c *connectionKeeper) writePump(ctx contexts.PromiseContext) {
	for {
		select {
		case ok := <-ctx.PromiseDone():
			//	drain the queue
		lp:
			for {
				select {
				case message := <-c.outputQueue:
					err := c.cnn.WebsocketWriter(ctx, message.Type(), message.Data())
					if err != nil {
						message.Sender().OnSendError(c, errs.New(err))
						continue
					}
				default:
					break lp
				}
			}
			ok()
			return
		case message := <-c.outputQueue:
			fmt.Println("c.outputQueue: ", len(c.outputQueue))
			err := c.cnn.WebsocketWriter(ctx, message.Type(), message.Data())
			if err != nil {
				message.Sender().OnSendError(c, errs.New(err))
				continue
			}
		}
	}
}

func (c *connectionKeeper) dispatchPump(ctx contexts.ContextGroup) {
	for {
		select {
		case ok := <-ctx.PromiseDone():
			//	drain the queue
		lp:
			for {
				select {
				case message := <-c.incomeQueue:
					c.cnn.Dispatch(ctx, message)
				default:
					break lp
				}
			}
			ok()
			return
		case message := <-c.incomeQueue:
			c.cnn.Dispatch(ctx, message)
		}
	}
}
