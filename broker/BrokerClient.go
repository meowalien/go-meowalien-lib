package broker

import "fmt"

// Filter will filter the messages input and return true if the message you want to pickup.
type Filter func(interface{}) bool

type Client struct {
	C      chan interface{}
	filter Filter
	uuid   string
}

func (b *Client) Filter() Filter {
	return b.filter
}

func (b *Client) UUID() string {
	return b.uuid
}



var ErrClientClosed = fmt.Errorf("the clinet has allready be closed")

func (b *Client) Close() error {
	if b.C == nil {
		return nil
	}
	select {
	case _ ,ok := <-b.C:
		if !ok{
			return ErrClientClosed //fmt.Errorf("the clinet has allready be closed")
		}else {
			close(b.C)
		}
	default:
		close(b.C)
	}
	return nil
}
