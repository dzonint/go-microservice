package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dzonint/go-microservice/data"
	"github.com/dzonint/go-microservice/handlers"
	"github.com/dzonint/go-microservice/rabbitmq"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	productLog := log.New()
	ph := handlers.NewProducts(productLog)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	getRouter.HandleFunc("/{id:[0-9]+}", ph.GetProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.RemoveProduct)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	data.InitDB("products.db")
	err := data.PopulateDB()
	if err != nil {
		log.Warning("Failed to populate database `products`")
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "-rabbitmq" {
			data.InitDB("users.db")

			uh := handlers.NewUsers()
			getRouter.HandleFunc("/users", uh.GetUsers)

			rmq := rabbitmq.NewRabbitMQService()
			go func() {
				rmq.GenerateConsumer()
			}()

			go func() {
				rmq.GenerateProducer()
			}()
		}
	}

	go func() {
		log.Fatal(s.ListenAndServe())
	}()
	log.Info("Server successfully started")

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Info("Received terminate, graceful shutdown ", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(timeoutContext)
}
