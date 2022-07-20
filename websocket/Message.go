package websocket

import "io"

type Message interface {
	io.Reader
}
