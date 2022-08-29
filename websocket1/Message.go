package websocket1

import "io"

type Message interface {
	Type() MessageType
	io.Reader
}
type MessageType int

/*
https://www.rfc-editor.org/rfc/rfc6455#section-5.6

0x1 (Text), 0x2 (Binary)
*/
const (
	// MessageText is for UTF-8 encoded text messages like JSON.
	MessageText MessageType = iota + 1
	// MessageBinary is for binary messages like protobufs.
	MessageBinary
)

func (t MessageType) Valid() bool {
	switch t {
	case MessageText:
	case MessageBinary:
	default:
		return false
	}
	return true
}

type TextMessage interface {
	Message
}
type BinaryMessage interface {
	Message
}

func NewTextMessage(data []byte) TextMessage {

}

func NewBinaryMessage(data []byte) BinaryMessage {

}
