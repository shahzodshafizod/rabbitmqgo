package main

import (
	"context"
	"io"
	"net/http"

	"github.com/shahzodshafizod/rabbitmqgo/pkg/config"
	"github.com/shahzodshafizod/rabbitmqgo/pkg/rabbitmq"
	"go.uber.org/fx"
)

func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var (
		pub rabbitmq.Publisher
		sub rabbitmq.Subscriber
	)
	fx.New(
		fx.Invoke(config.Read),
		fx.Provide(rabbitmq.New),
		fx.Populate(&pub),
		fx.Populate(&sub),
	).Start(ctx)

	key := "THE_KEY"

	deliveries, err := sub.Subscribe(ctx, key)
	if err != nil {
		println("sub.Subscribe", err)
		return
	}

	go func() {
		for delivery := range deliveries {
			println("<- Received message", string(delivery.Body))
		}
	}()

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		if err := pub.Publish(ctx, key, data); err != nil {
			println("pub.Publish", err)
		}
	})
	http.ListenAndServe(":8080", nil)
}
