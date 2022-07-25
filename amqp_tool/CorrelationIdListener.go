package amqp_tool

import (
	"fmt"
	"sync"

	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/streadway/amqp"
)

func NewCorrelationIdListener(ch *amqp.Channel, queueName string) CorrelationIdListener {
	return &correlationIdListener{ch: ch, queueName: queueName}
}

type CorrelationIdListener interface {
	AddListener(key string, listener ListenerFunc) (err error)
	RemoveListener(key string)
	Start() (err error)
}

type StopListen func()

type ListenerFunc func(delivery amqp.Delivery, ls StopListen)

type correlationIdListener struct {
	listenerMap sync.Map
	ch          *amqp.Channel
	queueName   string
}

func (r *correlationIdListener) Start() (err error) {
	div, err := r.ch.Consume(r.queueName, "", false, false, false, false, nil)
	if err != nil {
		err = errs.New(err)
		return
	}
	go func() {
		for d := range div {
			delivery := d
			listener, ok := r.getListener(delivery.CorrelationId)
			if ok {
				listener(delivery, func() {
					fmt.Println("callback been called")
					r.RemoveListener(delivery.CorrelationId)
				})
				continue
			} else {
				fmt.Println("not subscribed delivery")
				err = delivery.Reject(false)
				if err != nil {
					fmt.Printf("error when Reject: %s\n", err.Error())
					//return
					continue
				}

				continue
			}
		}
		fmt.Println("Consumer close")
	}()
	return
}

func (r *correlationIdListener) AddListener(key string, listener ListenerFunc) (err error) {

	fmt.Println("AddListener: ", key)
	_, loaded := r.listenerMap.LoadOrStore(key, listener)
	if loaded {
		err = errs.New("the Listener on key %s is already exist", key)
		return
	}
	//r.listenerMap.Range(func(key, value interface{}) bool {
	//	fmt.Printf("AddListener - key:%s , val:%+v\n" , key , value)
	//	return true
	//})

	return
}
func (r *correlationIdListener) RemoveListener(key string) {
	fmt.Println("RemoveListener: ", key)
	r.listenerMap.Delete(key)
	//r.listenerMap.Range(func(key, value interface{}) bool {
	//	fmt.Printf("RemoveListener - key:%s , val:%+v\n" , key , value)
	//	return true
	//})
}

func (r *correlationIdListener) getListener(key string) (lf ListenerFunc, ok bool) {
	l, ok := r.listenerMap.Load(key)
	if !ok {
		return
	}
	lf = l.(ListenerFunc)
	return
}
