package controller

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/shafizod/rabbitmqgo/internal/rabbitmq"

	"github.com/julienschmidt/httprouter"
	"github.com/streadway/amqp"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
)

func Start(lifecycle fx.Lifecycle, rabbit rabbitmq.Service) {
	router := httprouter.New()
	router.POST("/publish", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		body := gjson.GetBytes(data, "body").String()
		if err := rabbit.Publish(
			"",    // exchange
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			}); err != nil {
			http.Error(w, "Failed to publish a message: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Congrats with sending message: %s", body)))
	})

	srv := &http.Server{
		Addr:         os.Getenv("APP_ADDRESS"),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 120 * time.Second,
		Handler:      router,
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(c context.Context) error {
			srv.Shutdown(c)
			return nil
		},
	})
}
