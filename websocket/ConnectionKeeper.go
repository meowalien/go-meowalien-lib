package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/meowalien/go-meowalien-lib/uuid"
)

type ConnectionKeeper interface {
	CloseConnection() error
	Open(writer http.ResponseWriter, request *http.Request, responseHeader http.Header) error
	SentBinaryMessage(message ...BinaryMessage)
	SentText(message ...string)
	SentMessage(message ...Message)
	SentJson(s ...interface{}) error
	SetHook(lifeTimeHook LifeTimeHook)
	UUID() string
	ConnectionClosed() bool
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type defaultLogger struct{}

func (d defaultLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (d defaultLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (d defaultLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (d defaultLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

const (
	DefaultReadBufferSize         = 1024
	DefaultWriteBufferSize        = 1024
	DefaultWriteWait              = 10 * time.Second
	DefaultPongWait               = 60 * time.Second
	DefaultPingPeriod             = (DefaultPongWait * 9) / 10
	DefaultMaxMessageSize         = 512
	DefaultSentMessageQueueLength = 256
)

var DefaultLogger = defaultLogger{}
var DefaultCheckOrigin = func(r *http.Request) bool {
	return true
}

type Option struct {
	Logger                 Logger
	ReadBufferSize         int
	WriteBufferSize        int
	SentMessageQueueLength int
	MaxMessageSize         int64
	PongWait               time.Duration
	PingPeriod             time.Duration
	WriteWait              time.Duration
	CheckOrigin            func(r *http.Request) bool
	// 禁止重複連線的唯一使識別符號
	ConnectionOwner string
	// default true
	AllowMultipleConnection bool
	// default false
	AutoDisconnectOldConnections bool

	ReadMiddleWare func(io.Reader)
}

func defaultOption() Option {
	return Option{
		ConnectionOwner:              "",
		AllowMultipleConnection:      true,
		AutoDisconnectOldConnections: false,
		CheckOrigin:                  DefaultCheckOrigin,
		ReadBufferSize:               DefaultReadBufferSize,
		WriteBufferSize:              DefaultWriteBufferSize,
		Logger:                       DefaultLogger,
		SentMessageQueueLength:       DefaultSentMessageQueueLength,
		MaxMessageSize:               DefaultMaxMessageSize,
		PongWait:                     DefaultPongWait,
		PingPeriod:                   DefaultPingPeriod,
		WriteWait:                    DefaultWriteWait,
	}
}

type OptionModifier func(option *Option)

func NewConnectionKeeper(optionModifier ...OptionModifier) (ck ConnectionKeeper) {
	option := defaultOption()
	for _, modifier := range optionModifier {
		modifier(&option)
	}

	ck = &connectionKeeper{
		uuid:   uuid.NewUUID("CK"),
		Option: option,
	}
	return ck
}

type connectionKeeper struct {
	Option
	conn               *websocket.Conn
	uuid               string
	sentMessageChannel chan Message
	LifeTimeHook
	//connectionClosed   bool
	readPumpClosed  bool
	writePumpClosed bool
}

func (c *connectionKeeper) UUID() string {
	return c.uuid
}

func (c *connectionKeeper) SetHook(lifeTimeHook LifeTimeHook) {
	c.LifeTimeHook = lifeTimeHook
}

func (c *connectionKeeper) SentBinaryMessage(message ...BinaryMessage) {
	for _, rawMessage := range message {
		c.SentMessage(rawMessage)
	}
}
func (c *connectionKeeper) SentMessage(message ...Message) {
	for _, rawMessage := range message {
		c.sentMessageChannel <- rawMessage
	}
}
func (c *connectionKeeper) SentText(message ...string) {
	for _, rawMessage := range message {
		c.SentMessage(NewTextMessage(rawMessage))
	}
}

func (c *connectionKeeper) SentJson(s ...interface{}) error {
	for _, i2 := range s {
		bf := bytes.NewBuffer([]byte{})
		jsonEncoder := json.NewEncoder(bf)
		jsonEncoder.SetEscapeHTML(false)
		err := jsonEncoder.Encode(i2)
		if err != nil {
			return err
		}
		c.SentText(bf.String())
	}
	return nil
}

const writeWait = time.Second

func (c *connectionKeeper) Open(writer http.ResponseWriter, request *http.Request, responseHeader http.Header) (err error) {
	err = c.multipleConnectionProcess()
	if err != nil {
		return err
	}
	c.sentMessageChannel = make(chan Message, c.SentMessageQueueLength)

	var websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  c.ReadBufferSize,
		WriteBufferSize: c.WriteBufferSize,
		CheckOrigin:     c.CheckOrigin,
	}
	c.conn, err = websocketUpgrader.Upgrade(writer, request, responseHeader)
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	c.conn.SetReadLimit(c.MaxMessageSize)
	err = c.conn.SetReadDeadline(time.Now().Add(c.PongWait))
	if err != nil {
		return fmt.Errorf("error when SetReadDeadline: %s", err.Error())
	}

	c.conn.SetPingHandler(func(message string) error {
		err := c.conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(writeWait))
		if err != nil {
			if err == websocket.ErrCloseSent {
				return nil
			} else if e, ok := err.(net.Error); ok && e.Temporary() {
				return nil
			} else {
				return err
			}
		}

		return c.OnPing(message)
	})

	c.conn.SetPongHandler(func(msg string) error {
		e := c.conn.SetReadDeadline(time.Now().Add(c.PongWait))
		if e != nil {
			return e
		}
		//if c.lifeTimeHook.OnPong != nil {
		return c.OnPong(msg)
		//}
		//return nil
	})

	go c.readPump()
	go c.writePump()

	if c.LifeTimeHook == nil {
		c.LifeTimeHook = emptyDispatcher{}
	}

	c.OnOpenConnection(c.uuid)

	return nil
}

func (c *connectionKeeper) readPump() {
	for {
		messageType, reader, e := c.conn.NextReader()
		if e != nil {
			// 前端關閉的住況
			if closeError, ok := e.(*websocket.CloseError); ok {
				c.Logger.Warnf("Connection %s readPump close, type: %s", c.uuid, WebsocketCloseCodeNumberToString(closeError.Code))
			} else {
				c.Logger.Errorf("NextReader error: %s\n", e.Error())
			}
			c.readPumpClosed = true
			err := c.CloseConnection()
			if err != nil {
				c.Logger.Errorf("error when CloseConnection in readPump: ", err.Error())
			}
			return
		}

		if c.Option.ReadMiddleWare != nil {
			c.Option.ReadMiddleWare(reader)
		}

		//if conf.DEBUG_MOD {
		//	reader = log.ReaderLogger("WS_REQUEST", reader)
		//}

		message := &binaryMessage{
			Reader: reader,
		}

		switch messageType {
		//default:
		//	switch messageType {
		case websocket.TextMessage:
			c.OnTextMessage(&textMessage{
				message,
			})
		case websocket.BinaryMessage:
			c.OnBinaryMessage(message)
		default:
			c.Logger.Errorf("Unknown event: %v", message)
			//}
		}
	}
}

const MaxPingFailCount int = 10

func (c *connectionKeeper) writePump() {
	pingTimer := time.NewTimer(c.PingPeriod)
	var pingFailCount int
loop:
	for {
		select {
		case message, ok := <-c.sentMessageChannel:
			if !ok {
				break loop
			}
			err := c.conn.SetWriteDeadline(time.Now().Add(c.WriteWait))
			if err != nil {
				c.Logger.Errorf(errs.WithLine(err).Error())
				continue
			}
			var writer io.WriteCloser
			switch t := message.(type) {
			case TextMessage:
				writer, err = c.conn.NextWriter(websocket.TextMessage)
			case BinaryMessage:
				writer, err = c.conn.NextWriter(websocket.BinaryMessage)
			default:
				c.Logger.Errorf("not supported message type: %T", t)
				continue

			}

			if err != nil {
				c.Logger.Errorf(errs.WithLine(err).Error())
				continue
			}

			_, err = io.Copy(writer, message)
			if err != nil {
				c.Logger.Errorf(errs.WithLine(err).Error())
				continue
			}

			if err = writer.Close(); err != nil {
				c.Logger.Errorf(errs.WithLine(err).Error())
				continue
			}
		case <-pingTimer.C:
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				c.Logger.Errorf("error when WriteMessage: %s", err.Error())
				pingFailCount++
				if pingFailCount < MaxPingFailCount {
					c.Logger.Debugf("try to ping again, count: %d", pingFailCount)
					pingTimer.Reset(c.PingPeriod)
					continue
				} else {
					c.Logger.Debugf("ping reach retry limit: %d, closing connection", MaxPingFailCount)
					break loop
				}
			}
			if pingFailCount != 0 {
				pingFailCount = 0
			}
			err = c.conn.SetWriteDeadline(time.Now().Add(c.WriteWait))
			if err != nil {
				c.Logger.Errorf("error when SetWriteDeadline: %s", err.Error())
				continue
			}
			pingTimer.Reset(c.PingPeriod)
		}
	}
	if !pingTimer.Stop() {
		select {
		case <-pingTimer.C:
		default:
		}
	}
	err := c.CloseConnection()
	if err != nil {
		c.Logger.Errorf("error when CloseConnection in writePump: ", err.Error())
	}
}

const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

//RFC_6455
func WebsocketCloseCodeNumberToString(errorCode int) string {
	switch errorCode {
	case CloseNormalClosure: //1000
		return "CloseNormalClosure"
	case CloseGoingAway: //1001
		return "CloseGoingAway"
	case CloseProtocolError: //1002
		return "CloseProtocolError"
	case CloseUnsupportedData: //1003
		return "CloseUnsupportedData"
	case CloseNoStatusReceived: //1005
		return "CloseNoStatusReceived"
	case CloseAbnormalClosure: //1006
		return "CloseAbnormalClosure"
	case CloseInvalidFramePayloadData: //1007
		return "CloseInvalidFramePayloadData"
	case ClosePolicyViolation: //1008
		return "ClosePolicyViolation"
	case CloseMessageTooBig: //1009
		return "CloseMessageTooBig"
	case CloseMandatoryExtension: //1010
		return "CloseMandatoryExtension"
	case CloseInternalServerErr: //1011
		return "CloseInternalServerErr"
	case CloseServiceRestart: //1012
		return "CloseServiceRestart"
	case CloseTryAgainLater: //1013
		return "CloseTryAgainLater"
	case CloseTLSHandshake: //1015
		return "CloseTLSHandshake"
	default:
		return "Unknown"
	}
}

func (c *connectionKeeper) CloseConnection() error {
	if c.ConnectionClosed() {
		return nil
	}

	c.Logger.Infof("closing %s websocket Connection... \n", c.uuid)

	if !c.writePumpClosed {
		c.writePumpClosed = true
		// 關閉寫入通道
		close(c.sentMessageChannel)
	}

	if !c.readPumpClosed {
		c.readPumpClosed = true
		err := c.conn.Close()
		if err != nil {
			return errs.WithLine(err)
		}
	}

	deleteConnectionWithTargetUUID(c.ConnectionOwner)

	//if c.LifeTimeHook != nil {
	c.OnCloseConnection(c.UUID())
	//}
	return nil
}

var ErrMultipleConnection = fmt.Errorf("tring to make ErrMultipleConnection on same ConnectionOwner")

// 處理多重連線的狀況
func (c *connectionKeeper) multipleConnectionProcess() error {
	if c.AllowMultipleConnection {
		return nil
	}

	conn, exist := getConnectionOnTarget(c.ConnectionOwner)
	if !exist {
		return nil
	}

	if !c.AutoDisconnectOldConnections {
		return ErrMultipleConnection
	}

	err := conn.CloseConnection()
	if err != nil {
		return errs.WithLine(err)
	}
	cacheConnectionWithTargetUUID(c.ConnectionOwner, conn)

	return nil
}

func (c *connectionKeeper) ConnectionClosed() bool {
	return c.readPumpClosed && c.writePumpClosed
}

var targetConnectionMap = map[string]ConnectionKeeper{}

func cacheConnectionWithTargetUUID(targetUUID string, connectionUUID ConnectionKeeper) {
	targetConnectionMap[targetUUID] = connectionUUID
}
func deleteConnectionWithTargetUUID(targetUUID string) {
	delete(targetConnectionMap, targetUUID)
}

func getConnectionOnTarget(targetUUID string) (ConnectionKeeper, bool) {
	c, e := targetConnectionMap[targetUUID]
	return c, e
}

type emptyDispatcher struct {
}

func (e emptyDispatcher) OnOpenConnection(connectionID string) {
	return
}

func (e emptyDispatcher) OnCloseConnection(connectionID string) {
	return
}

func (e emptyDispatcher) OnTextMessage(message TextMessage) {
	return
}

func (e emptyDispatcher) OnBinaryMessage(message BinaryMessage) {
	return
}

func (e emptyDispatcher) OnPong(message string) error {
	return nil
}

func (e emptyDispatcher) OnPing(message string) error {
	return nil
}
