package main

import (
	"github.com/shafizod/rabbitmqgo/cmd/subscribe/controller"
	"github.com/shafizod/rabbitmqgo/internal/config"
	"github.com/shafizod/rabbitmqgo/internal/rabbitmq"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.Init),
		fx.Provide(rabbitmq.New),
		fx.Invoke(controller.Start),
	).Run()
}
