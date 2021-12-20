package broker

import (
	"core1/src/pkg/meowalien_lib/uuid"
	"fmt"
	"log"
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
	Unsubscribe(msgCh *Client)
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

var allBroker []*broker

func StopAllBroker() {
	for _, b := range allBroker {
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
	allBroker = append(allBroker, b)
	for {
		select {
		case <-b.stopCh:
			for msgCh := range subs {
				err := msgCh.Close()
				if err != nil {
					fmt.Printf("error when close %s Client: %s\n", msgCh.UUID(), err.Error())
				}
			}
			return
		case msgCh := <-b.subscribeChan:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unSubscribeChan:
			err := msgCh.Close()
			if err != nil {
				log.Println("error when close Client: %s", err.Error())
			}
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
						bk.C <- msg
					}
				}
				threadLimiter <- struct{}{}
				doTransfer(msgCh)
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
	msgCh := &Client{C: make(chan interface{}, b.clientQueueSize), filter: filter, uuid: getUUID()}
	b.subscribeChan <- msgCh
	return msgCh
}

// Unsubscribe will make broker stop sending new message to the given Client cnd close the c channel.
func (b *broker) Unsubscribe(msgCh *Client) {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	b.unSubscribeChan <- msgCh
}

// Publish will broadcast the message to all subscribed Client.
func (b *broker) Publish(msg interface{}, except ...*Client) {
	if !b.isActive {
		panic("the broker is not active, please start it up.")
	}
	b.publishChan <- [2]interface{}{msg, except}
}
