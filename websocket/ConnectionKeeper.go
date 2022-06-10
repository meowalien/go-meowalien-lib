package websocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/meowalien/go-meowalien-lib/uuid"
)

const (
	DefaultPingWait               = time.Second *3
	DefaultWriteWait              = time.Second *3
	DefaultPongWait               = time.Second *3
	DefaultPingPeriod             = (DefaultPongWait * 9) / 10
	DefaultMaxPingFailCount       = 10
	DefaultReadBufferSize         = 1024
	DefaultWriteBufferSize        = 1024
	DefaultMaxMessageSize         = 512
	DefaultSentMessageQueueLength = 256
)

var DefaultLogger = defaultLogger{}
var DefaultCheckOrigin = func(r *http.Request) bool {
	return true
}
var ErrMultipleConnection = fmt.Errorf("tring to make ErrMultipleConnection on same ConnectionOwner")

type ConnectionKeeper interface {
	CloseConnection() error
	Open(writer http.ResponseWriter, request *http.Request, responseHeader http.Header) error
	SentBinaryMessage(message ...BinaryMessage)
	SentText(message ...string)
	SentMessage(message ...Message)
	SentJson(s ...interface{}) error
	SetHook(lifeTimeHook *LifeTimeHook)
	UUID() string
	ConnectionClosed() bool
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Option struct {
	// wait for ping send success
	PingWait                     time.Duration
	Logger                       Logger
	ReadBufferSize               int
	WriteBufferSize              int
	SentMessageQueueLength       int
	MaxMessageSize               int64
	PongWait                     time.Duration
	PingPeriod                   time.Duration
	WriteWait                    time.Duration
	CheckOrigin                  func(r *http.Request) bool
	ConnectionOwner              string
	AllowMultipleConnection      bool
	AutoDisconnectOldConnections bool
	ReadMiddleWare               func(io.Reader)
	MaxPingFailCount             int
}

type OptionModifier func(option *Option)

func defaultOption() Option {
	return Option{
		PingWait:                     DefaultPingWait,
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
		MaxPingFailCount:             DefaultMaxPingFailCount,
	}
}

func NewConnectionKeeper(optionModifier ...OptionModifier) (ck ConnectionKeeper) {
	option := defaultOption()
	for _, modifier := range optionModifier {
		modifier(&option)
	}

	ck = &connectionKeeper{
		websocketUpgrader: websocket.Upgrader{
			ReadBufferSize:  option.ReadBufferSize,
			WriteBufferSize: option.WriteBufferSize,
			CheckOrigin:     option.CheckOrigin,
		},
		sentMessageChannel:  make(chan Message, option.SentMessageQueueLength),
		targetConnectionMap: map[string]ConnectionKeeper{},
		uuid:                uuid.NewUUID("CK"),
		Option:              option,
	}
	return ck
}

type connectionKeeper struct {
	lifeTimeHook *LifeTimeHook
	Option
	targetConnectionMap map[string]ConnectionKeeper
	conn                *websocket.Conn
	uuid                string
	sentMessageChannel  chan Message
	readPumpClosed      bool
	writePumpClosed     bool
	websocketUpgrader   websocket.Upgrader
}

func (c *connectionKeeper) UUID() string {
	return c.uuid
}

func (c *connectionKeeper) SetHook(lifeTimeHook *LifeTimeHook) {
	c.lifeTimeHook = lifeTimeHook
}

func (c *connectionKeeper) SentBinaryMessage(message ...BinaryMessage) {
	for _, rawMessage := range message {
		c.SentMessage(rawMessage)
	}
}

func (c *connectionKeeper) SentMessage(message ...Message) {
	for _, rawMessage := range message {
		select {
		case c.sentMessageChannel <- rawMessage:
		default:
			fmt.Println("sentMessageChannel full, drop: ", message)
		}
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

func (c *connectionKeeper) Open(writer http.ResponseWriter, request *http.Request, responseHeader http.Header) (err error) {
	err = c.multipleConnectionProcess()
	if err != nil {
		return err
	}

	c.conn, err = c.websocketUpgrader.Upgrade(writer, request, responseHeader)
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	c.conn.SetReadLimit(c.MaxMessageSize)
	err = c.conn.SetReadDeadline(time.Now().Add(c.PongWait))
	if err != nil {
		return fmt.Errorf("error when SetReadDeadline: %w", err)
	}
	err = c.conn.SetWriteDeadline(time.Now().Add(c.WriteWait))
	if err != nil {
		return fmt.Errorf("error when SetReadDeadline: %w", err)
	}
	go c.readPump()
	go c.writePump()
	if c.lifeTimeHook != nil && c.lifeTimeHook.OnOpenConnection != nil {
		c.lifeTimeHook.OnOpenConnection(c.uuid)
	}
	return nil
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

	c.deleteConnectionWithTargetUUID(c.ConnectionOwner)
	if c.lifeTimeHook != nil && c.lifeTimeHook.OnCloseConnection != nil {
		c.lifeTimeHook.OnCloseConnection(c.UUID())
	}
	return nil
}

func (c *connectionKeeper) ConnectionClosed() bool {
	return c.readPumpClosed && c.writePumpClosed
}

func (c *connectionKeeper) readPump() {
	c.conn.SetPingHandler(func(message string) error {
		err := c.conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(c.WriteWait))
		if err != nil {
			if errors.Is(err, websocket.ErrCloseSent) {
				return nil
			}
			switch e := err.(type) {
			case net.Error:
				if e.Temporary() {
					return nil
				} else {
					return err
				}
			default:
				return err
			}
		}
		if c.lifeTimeHook != nil && c.lifeTimeHook.OnPing != nil {
			c.lifeTimeHook.OnPing(message)
		}
		return nil
	})

	c.conn.SetPongHandler(func(msg string) error {
		e := c.conn.SetReadDeadline(time.Now().Add(c.PongWait))
		if e != nil {
			return fmt.Errorf("error when SetReadDeadline: %w", e)
		}
		if c.lifeTimeHook != nil && c.lifeTimeHook.OnPong != nil {
			c.lifeTimeHook.OnPong(msg)
		}
		return nil
	})
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

		message := &binaryMessage{
			Reader: reader,
		}

		switch messageType {
		case websocket.TextMessage:
			if c.lifeTimeHook != nil && c.lifeTimeHook.OnTextMessage != nil {
				c.lifeTimeHook.OnTextMessage(&textMessage{
					message,
				})
			}

		case websocket.BinaryMessage:
			if c.lifeTimeHook != nil && c.lifeTimeHook.OnBinaryMessage != nil {
				c.lifeTimeHook.OnBinaryMessage(message)
			}
		default:
			c.Logger.Errorf("Unknown event: %v", message)
		}
	}
}

func (c *connectionKeeper) writePump() {
	pingTimer := time.NewTimer(c.PingPeriod)
	var pingCount int
	var pingFailCount int
loop:
	for {
		select {
		case message, ok := <-c.sentMessageChannel:
			if !ok {
				break loop
			}
			var err error
			var writer io.WriteCloser
			switch t := message.(type) {
			case TextMessage:
				writer, err = c.conn.NextWriter(websocket.TextMessage)
			case BinaryMessage:
				writer, err = c.conn.NextWriter(websocket.BinaryMessage)
			default:
				c.Logger.Errorf("not supported message type: %T\n", t)
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
			err = c.conn.SetWriteDeadline(time.Now().Add(c.WriteWait))
			if err != nil {
				c.Logger.Errorf(errs.WithLine(err).Error())
				continue
			}
		case <-pingTimer.C:
			pingCount++
			err := c.conn.WriteControl(websocket.PingMessage, []byte(strconv.Itoa(pingCount)), time.Now().Add(c.PingWait))
			if err != nil {
				c.Logger.Errorf("error when SentPing: %s", err.Error())
				pingFailCount++
				if pingFailCount < c.MaxPingFailCount {
					c.Logger.Debugf("try to ping again after %s, count: %d", c.PingPeriod.String(), pingFailCount)
					pingTimer.Reset(c.PingPeriod)
					continue
				} else {
					c.Logger.Debugf("ping reach retry limit: %d, closing connection", c.MaxPingFailCount)
					break loop
				}
			}
			pingFailCount = 0
			err = c.conn.SetWriteDeadline(time.Now().Add(c.WriteWait))
			if err != nil {
				c.Logger.Errorf("error when SetWriteDeadline: %s", err.Error())
				break loop
			}
			pingTimer.Reset(c.PingPeriod)
			if c.lifeTimeHook != nil && c.lifeTimeHook.AfterPingEmit != nil {
				c.lifeTimeHook.AfterPingEmit(pingCount)
			}
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

// 處理多重連線的狀況
func (c *connectionKeeper) multipleConnectionProcess() error {
	if c.AllowMultipleConnection {
		return nil
	}

	conn, exist := c.getConnectionOnTarget(c.ConnectionOwner)
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
	c.cacheConnectionWithTargetUUID(c.ConnectionOwner, conn)

	return nil
}

func (c *connectionKeeper) cacheConnectionWithTargetUUID(targetUUID string, connectionUUID ConnectionKeeper) {
	c.targetConnectionMap[targetUUID] = connectionUUID
}

func (c *connectionKeeper) deleteConnectionWithTargetUUID(targetUUID string) {
	delete(c.targetConnectionMap, targetUUID)
}

func (c *connectionKeeper) getConnectionOnTarget(targetUUID string) (ConnectionKeeper, bool) {
	ck, e := c.targetConnectionMap[targetUUID]
	return ck, e
}
