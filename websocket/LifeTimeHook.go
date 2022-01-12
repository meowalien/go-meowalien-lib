package websocket

type LifeTimeHook interface {
	OnOpenConnection(connectionID string)
	OnCloseConnection(connectionID string)
	OnTextMessage(message TextMessage)
	OnBinaryMessage(message BinaryMessage)
	OnPong(message BinaryMessage)
	OnPing(message BinaryMessage)
}
