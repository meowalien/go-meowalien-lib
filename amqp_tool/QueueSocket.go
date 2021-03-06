package amqp_tool

import (
	"fmt"
	"sync/atomic"

	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/streadway/amqp"
)

type QueueSocket interface {
	PushTopic(topic string, skd SocketData, callback ListenerFunc) (err error)
}

type SocketData struct {
	CorrelationId string
	ContentType   string
	Body          []byte
}

func NewQueueSocket(channel *amqp.Channel, exchangeName string) (sk QueueSocket, err error) {
	err = channel.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		err = errs.New(err)
		return
	}
	// 回應的queue
	responseQueue, err := channel.QueueDeclare(
		"",
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		err = errs.New(err)
		return
	}

	responseListener := NewCorrelationIdListener(channel, responseQueue.Name)
	err = responseListener.Start()
	if err != nil {
		err = errs.New(err)
		return
	}
	sk = &socketKeeper{listener: responseListener, channel: channel, queueName: responseQueue.Name, exchangeName: exchangeName}
	return
}

type socketKeeper struct {
	listener     CorrelationIdListener
	channel      *amqp.Channel
	queueName    string
	exchangeName string
}

var count uint64

func newID(s string) string {
	n := atomic.AddUint64(&count, 1)
	return fmt.Sprintf("%d%s", n, s)
}

func (w *socketKeeper) PushTopic(topic string, skd SocketData, callback ListenerFunc) (err error) {
	//debug.PrintStack()
	fmt.Println("PushTopic: ", topic)
	id := newID(skd.CorrelationId)
	if callback != nil {
		err = w.listener.AddListener(id, callback)
		if err != nil {
			err = errs.New(err)
			return
		}
	}

	err = w.channel.Publish(
		w.exchangeName, // exchange
		topic,          // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:   skd.ContentType,
			CorrelationId: id,
			ReplyTo:       w.queueName,
			Body:          skd.Body,
		})
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
