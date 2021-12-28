package broker

import (
	"context"
	"fmt"
)

// Filter will filter the messages input and return true if the message you want to pickup.
type Filter func(interface{}) bool

type Client struct {
	c      chan interface{}
	filter Filter
	uuid       string
	onNewEvent func(newEvent interface{})
	onClose    func(client *Client)
	onError func(err error)
}

func (b *Client) Filter() Filter {
	return b.filter
}

func (b *Client) UUID() string {
	return b.uuid
}

var ErrClientClosed = fmt.Errorf("the clinet has allready be closed")

func (b *Client) Close() error {
	if b.c == nil {
		return nil
	}
	select {
	case _, ok := <-b.c:
		if !ok {
			return ErrClientClosed //fmt.Errorf("the clinet has allready be closed")
		} else {
			close(b.c)
		}
	default:
		close(b.c)
	}
	return nil
}
func (b *Client) Listen(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				err := b.Close()
				if err != nil && b.onError != nil {
					b.onError(err)
				}
				return
			case ev, ok := <-b.c:
				if !ok {
					if b.onClose != nil {
						b.onClose(b)
					}
					return
				}
				if b.onNewEvent != nil {
					b.onNewEvent(ev)
				}
			}
		}
	}()
}

func (b *Client) OnError(f func(err error)) {
	b.onError = f
}
func (b *Client) OnNewEvent(f func(ev interface{})) {
	b.onNewEvent = f
}
func (b *Client) OnClose(f func(ev *Client)) {
	b.onClose = f
}
