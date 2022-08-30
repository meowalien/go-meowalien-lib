package websocket1

import (
	"bytes"
	"context"
	"fmt"
)

func NewBinaryMessage(sender MessageSender, data []byte) BinaryMessage {
	return &binaryMessage{
		msgSender: sender,
		sender:    donothingSender,
		data:      data,
		Reader:    bytes.NewReader(data),
	}
}

type BinaryMessage interface {
	Message
}

type binaryMessage struct {
	data []byte
	*bytes.Reader
	sender    Sender
	msgSender MessageSender
}

func (b *binaryMessage) ReplyText(ctx context.Context, text string) (err error) {
	return b.msgSender.SendMessage(ctx, NewTextMessage(b.msgSender, text))

}

func (b *binaryMessage) ReplyBinary(ctx context.Context, text string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (b *binaryMessage) String() string {
	return fmt.Sprint(b.data)
}

func (b *binaryMessage) Type() MessageType {
	return MessageTypeBinary
}

func (b *binaryMessage) Data() []byte {
	return b.data
}

func (b *binaryMessage) Sender() Sender {
	return b.sender
}
