package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/karimLa/microservices/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProduct(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.Use(ph.MiddlewareValidateProduct)
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.Use(ph.MiddlewareValidateProduct)
	postRouter.HandleFunc("/", ph.AddProduct)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block until a signal is received.
	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
