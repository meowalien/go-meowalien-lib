package websocket1

func NewTextMessage(sender MessageSender, data string) TextMessage {
	return &textMessage{
		BinaryMessage: NewBinaryMessage(sender, []byte(data)),
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
