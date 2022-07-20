package websocket

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

type BinaryMessage interface {
	Message
	Binary() ([]byte, error)
	UnmarshalJson(a interface{}) error
}

type binaryMessage struct {
	io.Reader
	binaryCatch []byte
}

func (m *binaryMessage) Binary() ([]byte, error) {
	if m.binaryCatch != nil {
		return m.binaryCatch, nil
	}
	msgBytes, err := ioutil.ReadAll(m)
	if err != nil {
		return nil, err
	}

	m.binaryCatch = msgBytes
	m.Reader = bytes.NewReader(msgBytes)

	return msgBytes, nil
}

func (m binaryMessage) UnmarshalJson(a interface{}) error {
	bin, err := m.Binary()
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, a)
	if err != nil {
		return err
	}
	return nil
}
