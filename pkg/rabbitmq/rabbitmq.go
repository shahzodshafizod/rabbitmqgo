package rabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

type connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type Connection interface {
	Close() error
	Publish(exchange string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

type Params struct {
	Url       string
	QueueName string
}

func New(pms *Params) (Connection, error) {
	conn, err := amqp.Dial(pms.Url)
	if err != nil {
		return nil, errors.New("Failed to connect to RabbitMQ: " + err.Error())
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, errors.New("Failed to open a channel: " + err.Error())
	}
	queue, err := channel.QueueDeclare(
		pms.QueueName, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return nil, errors.New("Failed to declare a queue: " + err.Error())
	}

	return &connection{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

func (c *connection) Close() error {
	if err := c.channel.Close(); err != nil {
		return errors.New("c.channel.Close ERROR: " + err.Error())
	}
	if err := c.conn.Close(); err != nil {
		return errors.New("c.conn.Close ERROR: " + err.Error())
	}
	return nil
}

func (c *connection) Publish(exchange string, mandatory, immediate bool, msg amqp.Publishing) error {
	return c.channel.Publish(exchange, c.queue.Name, mandatory, immediate, msg)
}

func (c *connection) Consume(consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(c.queue.Name, consumer, autoAck, exclusive, noLocal, noWait, args)
}
