package controller

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/shafizod/rabbitmqgo/internal/rabbitmq"
	"github.com/tidwall/gjson"

	"github.com/julienschmidt/httprouter"
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
		if !gjson.ValidBytes(data) {
			http.Error(w, "Invalid json: "+string(data), http.StatusBadRequest)
			return
		}
		if err := rabbit.Publish("merchants", data); err != nil {
			http.Error(w, "Failed to publish a message: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Congrats with sending message: %s", []byte(data))))
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
