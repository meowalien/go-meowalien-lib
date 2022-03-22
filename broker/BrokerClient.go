package broker

import "fmt"

//type SweetHelperLineBotLib interface {
//	Close() error
//	UUID() string
//	C() chan interface{}
//	Filter() Filter
//}

// Filter will filter the messages input and return true if the message you want to pickup.
type Filter func(interface{}) bool

type Client struct {
	// New messages will be received through c
	C      chan interface{}
	broker Broker
	filter Filter
	uuid   string
}

func (b *Client) Publish(msg interface{}) {
	b.broker.Publish(msg , b)
}

func (b *Client) Filter() Filter {
	return b.filter
}
func (b *Client) UUID() string {
	return b.uuid
}

var ErrClientClosed = fmt.Errorf("the clinet has allready be closed")

func (b *Client) close() {
	if b.C == nil {
		return
	}
	select {
	case _, ok := <-b.C:
		if !ok {
			return //fmt.Errorf("the clinet has allready be closed")
		} else {
			close(b.C)
		}
	default:
		close(b.C)
	}
	return
}

func (b *Client) Close() {
	b.broker.unsubscribe(b)
}
