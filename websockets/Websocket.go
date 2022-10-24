package websockets

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

type ConnectionAdapter struct {
	OnError         OnErrorCallback
	Dispatcher      Dispatcher
	WebsocketReader Reader
	WebsocketWriter Writer
	Close           func() error
}
type WebsocketMessageQueuer interface {
	Start()
	Close() error
	CloseReader() error
	CloseWriter() error
	SendMessage(ctx context.Context, message Message) error
}

func NewConnectionKeeper(cnn ConnectionAdapter, configs ...Config) (keeper WebsocketMessageQueuer) {
	conf := defaultConfig()
	for c := range configs {
		configs[c](&conf)
	}
	onece := sync.Once{}
	// to prevent second call of Close()
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

type MessageSender func(ctx context.Context, message Message) (err error)
type Dispatcher func(ctx context.Context, msg Message)
type OnErrorCallback func(err error)
type Reader func(ctx context.Context) (msgType MessageType, data []byte, err error)
type Writer func(ctx context.Context, typ MessageType, p []byte) (err error)

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

type connectionKeeper struct {
	incomeQueue      chan Message
	outputQueue      chan Message
	cnn              ConnectionAdapter
	readTimeout      time.Duration
	pingInterval     time.Duration
	rootContextGroup contexts.ContextGroup
	dispatchPumpCtx  contexts.ContextGroup
	writePumpCtx     contexts.ContextGroup
	readPumpCtx      contexts.ContextGroup
}

func (c *connectionKeeper) CloseReader() (err error) {
	c.readPumpCtx.Close()
	return
}

func (c *connectionKeeper) CloseWriter() (err error) {
	c.readPumpCtx.Close()
	return
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
	c.writePumpCtx = contexts.NewContextGroup(c.rootContextGroup)
	c.dispatchPumpCtx = c.writePumpCtx.ChildGroup()
	c.readPumpCtx = c.dispatchPumpCtx.ChildGroup()

	go c.writePump(c.writePumpCtx)       // close 3st
	go c.dispatchPump(c.dispatchPumpCtx) // close 2st
	go c.readPump(c.readPumpCtx)         // close 1st
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
			fmt.Println("readPump done")
			// break when the context is done, but the dispatchPump is still running, it will drain the queue
			return
		default:
			msgtype, data, err := c.cnn.WebsocketReader(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
					err = c.Close()
					if err != nil {
						err = errs.New(err)
						c.cnn.OnError(err)
						return
					}
					return
				}
				c.cnn.OnError(errs.New("error when reading:%w , data:%v , message type: %v", err, data, msgtype))
				return
			}
			var message Message
			switch msgtype {
			case MessageTypeBinary:
				message = NewBinaryMessage(data)
			case MessageTypeText:
				message = NewTextMessage(string(data))
			default:
				panic("unknown message type")
			}
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
				c.cnn.OnError(errs.New("%w, dropping message: %v", ErrIncomeQueueFull, message))
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
					err := c.cnn.WebsocketWriter(ctx, message.Type(), message.Bytes())
					if err != nil {
						c.cnn.OnError(errs.New("error when writing:%w , data:%v , message type: %v", err, message.Bytes(), message.Type()))
						continue
					}
				default:
					break lp
				}
			}
			ok()
			fmt.Println("writePump done")
			return
		case message := <-c.outputQueue:
			//fmt.Println("c.outputQueue: ", len(c.outputQueue))
			err := c.cnn.WebsocketWriter(ctx, message.Type(), message.Bytes())
			if err != nil {
				c.cnn.OnError(errs.New("error when writing:%w , data:%v , message type: %v", err, message.Bytes(), message.Type()))
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
					c.cnn.Dispatcher(ctx, message)
				default:
					break lp
				}
			}
			ok()
			return
		case message := <-c.incomeQueue:
			c.cnn.Dispatcher(ctx, message)
		}
	}
}
