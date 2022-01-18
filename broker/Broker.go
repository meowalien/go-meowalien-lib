package broker

import (
	"github.com/meowalien/go-meowalien-lib/uuid"
	"log"
	"runtime/debug"
	"time"
)

const DefaultChanSize = 1
const DefaultSubscribeChanSize = 1
const DefaultUnSubscribeChanSize = 1
const DefaultClientQueueSize = 5
const MaximumBroadcastThreads = 1000

type Options struct {
	PublishChanSize     int
	SubscribeChanSize   int
	UnSubscribeChanSize int
	ClientQueueSize     int
}

var threadLimiter = make(chan struct{}, MaximumBroadcastThreads)

type Broker interface {
	//start()
	Subscribe(filter Filter) *Client
	unsubscribe(msgCh *Client)
	Publish(msg interface{}, except ...*Client)
}

// broker is a bridge between multiple Client, it will transfer data between them.
type broker struct {
	// signal to stop broker
	stopCh chan struct{}
	// sent message to all Client
	publishChan chan [2]interface{}
	// subscribe new Client
	subscribeChan chan *Client
	// Unsubscribe Client
	unSubscribeChan chan *Client

	isActive        bool
	clientQueueSize int
}

var activeBrokers []*broker

func StopAllBroker() {
	for _, b := range activeBrokers {
		b.stop()
	}
}

// NewBroker create a new broker according to given option, will create a default broker if the given option is nil
func NewBroker(option *Options) Broker {
	if option == nil {
		option = &Options{
			PublishChanSize:     DefaultChanSize,
			SubscribeChanSize:   DefaultSubscribeChanSize,
			UnSubscribeChanSize: DefaultUnSubscribeChanSize,
			ClientQueueSize:     DefaultClientQueueSize,
		}
	}
	bk := &broker{
		stopCh:          make(chan struct{}),
		publishChan:     make(chan [2]interface{}, option.PublishChanSize),
		subscribeChan:   make(chan *Client, option.SubscribeChanSize),
		unSubscribeChan: make(chan *Client, option.UnSubscribeChanSize),
		clientQueueSize: option.ClientQueueSize,
	}

	go bk.start()
	return bk
}

func getUUID() string {
	return uuid.NewUUID("B")
}

// start should be called before broker use, it starts up the broker
func (b *broker) start() {
	subs := map[*Client]struct{}{}
	b.isActive = true
	defer func() { b.isActive = false }()
	activeBrokers = append(activeBrokers, b)
lp:
	for {
		select {
		case <-b.stopCh:
			for msgCh := range subs {
				//err :=
				msgCh.close()
				//if err != nil {
				//	fmt.Printf("error when close %s Client: %s\n", msgCh.UUID(), err.Error())
				//}
			}
			break lp
		case msgCh := <-b.subscribeChan:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unSubscribeChan:
			//err :=
			msgCh.close()
			//if err != nil {
			//	log.Printf("error when close Client: %s\n", err.Error())
			//}
			delete(subs, msgCh)
		case m := <-b.publishChan:
			msg := m[0]
			allExcept := m[1].([]*Client)

			for msgCh := range subs {
				doTransfer := func(bk *Client) {
					if allExcept != nil {
						for _, exceptMsgCh := range allExcept {
							if exceptMsgCh == bk {
								return
							}
						}
					}
					if !b.isActive {
						return
					}

					if bk.Filter() != nil && bk.Filter()(msg) {
						if msgCh == nil || msgCh.C == nil {
							return
						}
						bk.C <- msg
					}
				}

				threadLimiter <- struct{}{}
				go doTransfer(msgCh)
				<-threadLimiter
			}
		}
	}
}

// stop will stop the broker
func (b *broker) stop() {
	b.isActive = false
	close(b.stopCh)
}

// Subscribe will create a new Client which Listen on new published message
func (b *broker) Subscribe(filter Filter) *Client {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	msgCh := &Client{broker: b, C: make(chan interface{}, b.clientQueueSize), filter: filter, uuid: getUUID()}
	b.subscribeChan <- msgCh
	return msgCh
}

// Unsubscribe will make broker stop sending new message to the given Client cnd close the c channel.
func (b *broker) unsubscribe(msgCh *Client) {
	//fmt.Println("Unsubscribe Client: ",msgCh)
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	if msgCh == nil {
		log.Println("the msgCh is nil")
		debug.PrintStack()
		return
	}
	b.unSubscribeChan <- msgCh
}

// Publish will broadcast the message to all subscribed Client.
func (b *broker) Publish(msg interface{}, except ...*Client) {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	timeoutTick := time.NewTimer(time.Second * 3)
	defer timeoutTick.Stop()
	select {
	case b.publishChan <- [2]interface{}{msg, except}:
		return
	case <-timeoutTick.C:
		log.Println("timeout when push to publishChan")
		debug.PrintStack()
		return
	}
}
