package websockets

import (
	"fmt"
	"io"
)

type Message interface {
	io.Reader
	fmt.Stringer
	Type() MessageType
	Data() []byte
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
