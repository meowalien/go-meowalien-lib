package websocket

import (
	"fmt"
	"io/ioutil"
	"net"
	http2 "net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestConnectionKeeper(t *testing.T) {
	t.Run("PingPong", TestConnectionKeeper_PingPong)
	//t.Run("OnCloseConnection", TestConnectionKeeper_OnCloseConnection)
}

func TestConnectionKeeper_PingPong(t *testing.T) {
	assert.NotPanics(t, func() {
		ck := NewConnectionKeeper(func(option *Option) {
			option.PingPeriod =
		})

		ck.SetHook(&LifeTimeHook{
			AfterPingEmit: func(count int) {
				fmt.Println("AfterPingEmit: ", count)
			},
			OnPong: func(message string) {
				fmt.Println("OnPong: ", message)
			},
			OnPing: func(message string) {
				fmt.Println("OnPing: ", message)
			},
		})

		wg := sync.WaitGroup{}

		c, _, err := makeWsClient(t, ck.Open)
		//c.SetReadDeadline(time.Now().Add(time.Second))
		c.SetPingHandler(func(appData string) error {
			fmt.Println("ping from server: ", appData)
			e := c.WriteControl(websocket.PongMessage , []byte(appData) ,time.Now().Add( time.Second ))
			assert.NoError(t, e)
			return nil
		})
		assert.NoError(t, err)
		defer c.Close()
		_, rd, err := c.NextReader()
		assert.NoError(t, err)
		fmt.Println("reader got")
		wg.Add(1)
		go func() {
			all, e := ioutil.ReadAll(rd)
			assert.NoError(t, e)
			fmt.Println("all: ", string(all))
			//assert.Equal(t, theMsg , string(all))
			wg.Done()
		}()
		wg.Wait()
	})
}

func makeWsClient(t *testing.T, open func(writer http2.ResponseWriter, request *http2.Request, responseHeader http2.Header) error) (wsconn *websocket.Conn, resp *http2.Response, err error) {
	mux := http2.NewServeMux()
	mux.HandleFunc("/", func(rw http2.ResponseWriter, r *http2.Request) {
		e := open(rw, r, nil)
		assert.NoError(t, e)
		return
	})
	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	httpServer := http2.Server{
		Addr:    listener.Addr().String(),
		Handler: mux,
	}

	go func() {
		err = httpServer.Serve(listener)
		if err != nil {
			assert.Error(t, err)
			return
		}
	}()

	dialer := websocket.DefaultDialer
	dialer.NetDial = net.Dial
	port := listener.Addr().String()
	addr := fmt.Sprintf("ws://127.0.0.1%s", port[strings.LastIndex(port, ":"):])
	wsconn, resp, err = dialer.Dial(addr, nil)
	return
}
