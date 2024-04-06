package rabbitmq

import (
	"github.com/streadway/amqp"
)

type Factory struct{}

func (f *Factory) NewConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (f *Factory) NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}
