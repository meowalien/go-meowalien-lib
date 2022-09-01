package websockets

import (
	"context"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"io"
	"net/http"
	"nhooyr.io/websocket"
	"sync"
	"testing"
	"time"
)

func TestWebsocket(t *testing.T) {
	wg := &sync.WaitGroup{}
	//wg.Add(1)
	go startServer(wg)
	time.Sleep(time.Second * 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, resp, err := websocket.Dial(ctx, "ws://localhost:9090", nil)
	if err != nil {
		err = errs.New(err)
		panic(err)
	}
	fmt.Println("resp: ", resp)

	keeper := NewConnectionKeeper(ConnectionAdapter{
		OnError: func(keeper ConnectionKeeper, err error) {
			fmt.Println("OnError: ", err)
		},
		Dispatcher: func(ctx context.Context, msg Message) {
			fmt.Println("Dispatch: ", msg)
		},
		WebsocketReader: func(ctx context.Context) (msgType MessageType, data []byte, err error) {
			msgTypeRaw, data, err := conn.Read(ctx)
			if err != nil {
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
					err = io.EOF
					return
				}
			}
			msgType = MessageType(msgTypeRaw)
			return
		},

		WebsocketWriter: func(ctx context.Context, typ MessageType, p []byte) (err error) {
			return conn.Write(ctx, websocket.MessageText, p)
		},
		Close: func() error {
			return conn.Close(websocket.StatusNormalClosure, "")
		},
	})
	go startClient(wg, keeper)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second * 1)
		for range ticker.C {
			err = keeper.SendMessage(context.Background(), NewTextMessage(keeper, "hello"))
			if err != nil {
				err = errs.New(err)
				panic(err)
			}
		}
	}()

	wg.Wait()
	fmt.Println("client done")
}

func startClient(wg *sync.WaitGroup, keeper ConnectionKeeper) {
	time.Sleep(time.Second * 1)
	defer wg.Done()

	keeper.Start()
}

func startServer(wg *sync.WaitGroup) {
	//defer wg.Done()
	hd := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{"*"},
		})
		if err != nil {
			err = errs.New(err)
			panic(err)
		}

		ap := newAdapter(conn)
		keeper := NewConnectionKeeper(ap)
		keeper.Start()

	})
	err := http.ListenAndServe(":9090", hd)
	if err != nil {
		return
	}
}

func newAdapter(conn *websocket.Conn) ConnectionAdapter {
	return ConnectionAdapter{
		OnError: func(keeper ConnectionKeeper, err error) {
			fmt.Println("OnError-server: ", err)
		},
		Dispatcher: func(ctx context.Context, msg Message) {
			fmt.Println("Dispatch-server: ", msg)
			err := msg.ReplyText(ctx, "hello-server")
			//err := keeper.SendMessage(ctx, NewTextMessage(keeper, "hello-server"))
			if err != nil {
				err = errs.New(err)
				panic(err)
			}
		},
		WebsocketReader: func(ctx context.Context) (msgType MessageType, data []byte, err error) {
			msgTypeRaw, data, err := conn.Read(ctx)
			if err != nil {
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
					err = io.EOF
					return
				}
			}
			msgType = MessageType(msgTypeRaw)
			return
		},

		WebsocketWriter: func(ctx context.Context, typ MessageType, p []byte) (err error) {
			return conn.Write(ctx, websocket.MessageText, p)
		},
		Close: func() error {
			return conn.Close(websocket.StatusNormalClosure, "")
		},
	}
}

func TestClosedFcChan(t *testing.T) {
	c := make(chan func(), 1)
	close(c)
	a := <-c
	a()
}
