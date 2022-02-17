package rabbitmq

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type connection struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	consumer string
}

type Connection interface {
	Close() error
	Publish(key string, body []byte) error
	Subscribe(queueName string, key string) (<-chan amqp.Delivery, error)
}

type Params struct {
	Url       string
	QueueName string
	Exchange  string
	Consumer  string
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
	if err := channel.ExchangeDeclare(
		pms.Exchange,       // name
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	); err != nil {
		return nil, errors.New("Failed to declare an exchange: " + err.Error())
	}
	return &connection{
		conn:     conn,
		channel:  channel,
		exchange: pms.Exchange,
		consumer: pms.Consumer,
	}, nil
}

func (c *connection) Close() error {
	if err := c.channel.Cancel(c.consumer, false); err != nil {
		return errors.New("c.channel.Cancel ERROR: " + err.Error())
	}
	if err := c.channel.Close(); err != nil {
		return errors.New("c.channel.Close ERROR: " + err.Error())
	}
	if err := c.conn.Close(); err != nil {
		return errors.New("c.conn.Close ERROR: " + err.Error())
	}
	return nil
}

func (c *connection) Publish(key string, body []byte) error {
	return c.channel.Publish(c.exchange, key, false, false, amqp.Publishing{
		ContentType: "text/json",
		Body:        body,
	})
}

func (c *connection) Subscribe(queueName string, key string) (<-chan amqp.Delivery, error) {
	queue, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.New("Failed to declare a queue: " + err.Error())
	}
	if err = c.channel.QueueBind(queue.Name, key, c.exchange, false, nil); err != nil {
		return nil, fmt.Errorf("bind queue %s: %w", queue.Name, err)
	}
	return c.channel.Consume(queueName, c.consumer, true, false, false, false, nil)
}
