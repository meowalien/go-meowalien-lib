package websocket

import (
	"context"
	"fmt"
	"io"
)

type Message interface {
	io.Reader
	fmt.Stringer
	Type() MessageType
	Data() []byte
	ReplyText(ctx context.Context, text string) (err error)
	ReplyBinary(ctx context.Context, text string) (err error)
	Sender() Sender
}
type Sender struct {
	OnSendError func(c ConnectionKeeper, err error)
}

var donothingSender = Sender{
	OnSendError: func(c ConnectionKeeper, err error) {},
}

type MessageType int

/*
https://www.rfc-editor.org/rfc/rfc6455#section-5.6

0x1 (Text), 0x2 (Binary)
*/
const (
	// MessageTypeText is for UTF-8 encoded text messages like JSON.
	MessageTypeText MessageType = iota + 1
	// MessageTypeBinary is for binary messages like protobufs.
	MessageTypeBinary
)

func (t MessageType) Valid() bool {
	switch t {
	case MessageTypeText:
	case MessageTypeBinary:
	default:
		return false
	}
	return true
}

func NewMessage(sender MessageSender, msgtype MessageType, data []byte) (message Message) {
	switch msgtype {
	case MessageTypeBinary:
		return NewBinaryMessage(sender, data)
	case MessageTypeText:
		return NewTextMessage(sender, string(data))
	default:
		panic("unknown message type")
	}
}
