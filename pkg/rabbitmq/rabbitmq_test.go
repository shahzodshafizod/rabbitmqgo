package rabbitmq

import (
	"context"
	"fmt"
	"testing"

	"github.com/shahzodshafizod/rabbitmqgo/pkg/config"
	"go.uber.org/fx"
)

/*
docker run -d --rm --name rabbitmq \
	-e RABBITMQ_DEFAULT_USER=user \
	-e RABBITMQ_DEFAULT_PASS=password \
	-p 5672:5672 rabbitmq
*/

// CONFIG_PATH_PREFIX=../../ go test -v -count=1 ./pkg/rabbitmq/ -run ^TestSub$
func TestSub(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var sub Subscriber
	fx.New(
		fx.Invoke(config.Read),
		fx.Provide(New),
		fx.Populate(&sub),
		fx.NopLogger,
	).Start(ctx)

	deliveries, err := sub.Subscribe(ctx, "THE_KEY")
	if err != nil {
		t.Fatal("sub.Subscribe", err)
	}

	for delivery := range deliveries {
		fmt.Println("<- Received message:", string(delivery.Body))
	}
}

// CONFIG_PATH_PREFIX=../../ go test -v -count=1 ./pkg/rabbitmq/ -run ^TestPub$
func TestPub(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var pub Publisher
	fx.New(
		fx.Invoke(config.Read),
		fx.Provide(New),
		fx.Populate(&pub),
		fx.NopLogger,
	).Start(ctx)

	if err := pub.Publish(ctx, "THE_KEY", []byte(`Hello World!`)); err != nil {
		t.Fatal("pub.Publish", err)
	}
}
