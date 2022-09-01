package websockets

import (
	"bytes"
	"fmt"
)

func NewBinaryMessage(data []byte) BinaryMessage {
	return &binaryMessage{
		sender: donothingSender,
		data:   data,
		Reader: bytes.NewReader(data),
	}
}

type BinaryMessage interface {
	Message
}

type binaryMessage struct {
	data []byte
	*bytes.Reader
	sender Sender
	raw    any
}

func (b *binaryMessage) Raw() any {
	return b.raw
}

//func (b *binaryMessage) Reply(ctx context.Context, msg Message) (err error) {
//	return b.msgSender(ctx, msg)
//}

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
