package websocket

type LifeTimeHook struct {
	OnOpenConnection  func(connectionID string)
	OnCloseConnection func(connectionID string)
	OnTextMessage     func(message TextMessage)
	OnBinaryMessage   func(message BinaryMessage)
	OnPong            func(message string)
	OnPing            func(message string)
	AfterPingEmit     func(count int)
}
