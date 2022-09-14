package websockets

import (
	"bytes"
	"fmt"
)

func NewBinaryMessage(data []byte) BinaryMessage {
	return &binaryMessage{
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
	raw any
}

func (b *binaryMessage) Raw() any {
	return b.raw
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
