package websocket1

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"io"
	"nhooyr.io/websocket"
	"testing"
)

type ConnectionAdapter struct {
	*websocket.Conn
}

func (c *ConnectionAdapter) AllocateMessage(msg Message) {
	decodeMsg, err := c.decodeMessage(msg)
	if err != nil {
		err = errs.New(err)
		panic(err)
		return
	}
	c.diliverDecodeResponse(decodeMsg)
}

type OperationCode = uint8 // 封包操作識別碼
type EventCode = uint8     // 封包事件識別碼

// DecodeResponse 解析成 code & payload
type DecodeResponse struct {
	OperationCode OperationCode
	EventCode     EventCode
	Data          []byte
}

func (c *ConnectionAdapter) decodeMessage(msg Message) (res *DecodeResponse, err error) {
	data, err := io.ReadAll(msg)
	dataLen := len(data)
	if dataLen < 3 {
		return nil, fmt.Errorf("data len=%v is invalid", dataLen)
	}

	operationCode := data[0]
	eventCode := binary.LittleEndian.Uint16(data[1:3])
	b := data[3:]

	res = &DecodeResponse{
		OperationCode: operationCode,
		EventCode:     uint8(eventCode),
		Data:          b,
	}

	return res, nil
}

func (c *ConnectionAdapter) Read(ctx context.Context) (msgType MessageType, data []byte, err error) {
	nhooyrMsgType, data, err := c.Conn.Read(context.TODO())
	msgType = MessageType(nhooyrMsgType)
	if !msgType.Valid() {
		err = errs.New("invalid message type: ", nhooyrMsgType)
		return
	}
	return
}

func (hdr *ConnectionAdapter) diliverDecodeResponse(msg *DecodeResponse) {
	payload, traceId, key, statusCode, err := hdr.getPayloadByCode(decodeMsg)

	// 如果收到回應，則移除對應請求
	if decodeMsg.OperationCode == operationcode.WebSocket_Response {
		key, _ = hdr.requestMap.Load(traceId)
		hdr.requestMap.Delete(traceId)
	}

	if err != nil {
		logger.Errorf("websocket/passMessageToRegister GetPayloadByCode err:%v", err)
		return
	}

	if len(key) == 0 {
		logger.Infof("websocket/passMessageToRegister got msg without key, op-code: %v, evt-code: %v, traceId: %v, payload: %v", decodeMsg.OperationCode, decodeMsg.EventCode, traceId, payload)
		return
	}

	target, exist := hdr.flowMap.Load(key)
	if !exist {
		// logger.Warnf("websocket/passMessageToRegister hdr.flowMap key=%v not found, traceId: %v, payload: %v", key, traceId, payload)
		return
	}

	if decodeMsg.OperationCode != operationcode.WebSocket_Ping {
		logger.Debug3f("websocket/passMessageToRegister, key: [%v], traceId: %v, op-code: %v, evt-code: %v, size: %v bytes", key, traceId, decodeMsg.OperationCode, decodeMsg.EventCode, len(decodeMsg.Data))
		logger.Debug2f("websocket/passMessageToRegister traceId: %v \n%v", traceId, string(decodeMsg.Data))
	}

	txName, ok := hdr.getTxNameByCode(key, decodeMsg.OperationCode, decodeMsg.EventCode)
	var tx *apm.Apm
	if ok {
		txType := util.SplitGuid(key)
		tx = apm.StartTransactionOptions(txName, txType, traceId)
	} else {
		// should log or not ?
	}
	target <- flow.NewApmRelayFlow(
		constant.ModuleId_WebSocket,
		constant.ModuleId_Game,
		decodeMsg.OperationCode,
		decodeMsg.EventCode,
		statusCode,
		traceId,
		key, // guid
		payload,
		tx,
	)
}

func TestWebsocketClient(t *testing.T) {
	ctx := context.Background()
	client, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
	if err != nil {
		err = errs.New(err)
		panic(err)
	}
	adapter := ConnectionAdapter{
		Conn: client,
	}

	keeper := NewConnectionKeeper(&adapter)
	keeper.Start(ctx)

}
