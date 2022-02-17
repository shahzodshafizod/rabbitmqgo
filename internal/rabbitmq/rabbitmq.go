package rabbitmq

import (
	"context"
	"errors"

	"github.com/shafizod/rabbitmqgo/pkg/rabbitmq"

	"github.com/streadway/amqp"
	"go.uber.org/fx"
)

type Service interface {
	Publish(key string, body []byte) error
	Subscribe(queueName string, key string) (<-chan amqp.Delivery, error)
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

func (s *service) Publish(key string, body []byte) error {
	return s.connection.Publish(key, body)
}
func (s *service) Subscribe(queueName string, key string) (<-chan amqp.Delivery, error) {
	return s.connection.Subscribe(queueName, key)
}
