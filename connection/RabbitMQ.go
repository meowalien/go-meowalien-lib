package connection

import (
	"fmt"
	"github.com/streadway/amqp"
)

type RabbitMQConfiguration struct {
	User     string
	Password string
	Host     string
	Port     string
}

func CreateRabbitMQConnection(c RabbitMQConfiguration) (conn *amqp.Connection, err error) {
	return  amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.User, c.Password, c.Host, c.Port))
}
