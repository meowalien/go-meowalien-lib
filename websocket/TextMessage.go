package websocket

import (
	"bytes"
	"github.com/meowalien/go-meowalien-lib/errs"
	"log"
)

func NewTextMessage(text string) TextMessage {
	return &textMessage{&binaryMessage{
		Reader: bytes.NewReader([]byte(text)),
	}}
}

type TextMessage interface {
	BinaryMessage
	Text() string
}

type textMessage struct {
	BinaryMessage
}

func (t textMessage) Text() string {
	bs, err := t.Binary()
	if err != nil {
		log.Println(errs.New(err))
		return ""
	}
	return string(bs)
}
