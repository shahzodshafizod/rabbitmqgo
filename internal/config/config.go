package config

import "github.com/shafizod/rabbitmqgo/pkg/rabbitmq"

func Init() (*rabbitmq.Params, error) {
	return &rabbitmq.Params{
		Url:       "amqp://guest:guest@rabbitmqgo:5672/",
		QueueName: "golang-queue",
	}, nil
}
