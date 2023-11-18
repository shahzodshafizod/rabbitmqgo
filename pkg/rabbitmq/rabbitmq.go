package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Publisher interface {
	Publish(ctx context.Context, key string, body []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, key string) (<-chan amqp.Delivery, error)
}

type client struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	exchange   string
	consumer   string
}

func New(lifecycle fx.Lifecycle) (Publisher, Subscriber, error) {

	var client = &client{
		queueName: viper.GetString("rabbitmq.queueName"),
		exchange:  viper.GetString("rabbitmq.exchange"),
		consumer:  viper.GetString("rabbitmq.consumer"),
	}

	var err error
	if client.connection, err = amqp.Dial(viper.GetString("rabbitmq.url")); err != nil {
		return nil, nil, err
	}

	if client.channel, err = client.connection.Channel(); err != nil {
		return nil, nil, err
	}

	if err = client.channel.ExchangeDeclare(
		client.exchange,    // name
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	); err != nil {
		return nil, nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			client.channel.Cancel(client.consumer, false)
			client.channel.Close()
			client.connection.Close()
			return nil
		},
	})

	return Publisher(client), Subscriber(client), nil
}

func (c *client) Publish(ctx context.Context, key string, body []byte) error {
	return c.channel.PublishWithContext(ctx, c.exchange, key, false, false, amqp.Publishing{
		ContentType: "text/json",
		Body:        body,
	})
}

func (c *client) Subscribe(ctx context.Context, key string) (<-chan amqp.Delivery, error) {
	queue, err := c.channel.QueueDeclare(
		c.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, nil
	}

	if err = c.channel.QueueBind(
		queue.Name,
		key,
		c.exchange,
		false,
		nil,
	); err != nil {
		return nil, nil
	}

	return c.channel.ConsumeWithContext(
		ctx,
		c.queueName,
		c.consumer,
		true,
		false,
		false,
		false,
		nil,
	)
}
