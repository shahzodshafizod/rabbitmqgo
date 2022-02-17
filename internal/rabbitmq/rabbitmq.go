package rabbitmq

import (
	"context"
	"errors"

	"github.com/shafizod/rabbitmqgo/pkg/rabbitmq"

	"github.com/streadway/amqp"
	"go.uber.org/fx"
)

type Service interface {
	Publish(exchange string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

type service struct {
	connection rabbitmq.Connection
}

func New(lifecycle fx.Lifecycle, pms *rabbitmq.Params) (Service, error) {
	connection, err := rabbitmq.New(pms)
	if err != nil {
		return nil, errors.New("rabbitmq.New ERROR: " + err.Error())
	}
	lifecycle.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			return connection.Close()
		},
	})
	return &service{connection: connection}, nil
}

func (s *service) Publish(exchange string, mandatory, immediate bool, msg amqp.Publishing) error {
	return s.connection.Publish(exchange, mandatory, immediate, msg)
}
func (s *service) Consume(consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return s.connection.Consume(consumer, autoAck, exclusive, noLocal, noWait, args)
}
