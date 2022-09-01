package websockets

func NewTextMessage(data string) TextMessage {
	return &textMessage{
		BinaryMessage: NewBinaryMessage([]byte(data)),
	}
}

type TextMessage interface {
	Message
}

type textMessage struct {
	BinaryMessage
}

func (t *textMessage) String() string {
	return string(t.Data())
}

func (t *textMessage) Type() MessageType {
	return MessageTypeText
}
