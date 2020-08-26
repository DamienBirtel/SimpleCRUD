package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/DamienBirtel/SimpleCRUD/handlers"
	"github.com/gorilla/mux"
)

var bindAddress string

func init() {
	bindAddress = os.Getenv("BINDADDRESS")
}

func main() {

	// create a new logger
	l := log.New(os.Stdout, "SimpleCRUD", log.LstdFlags)

	// create a new servemux
	m := mux.NewRouter()

	// register the handler
	m.HandleFunc("/", handlers.Get)

	// create a server
	s := http.Server{
		Addr:         bindAddress,
		Handler:      m,
		ErrorLog:     l,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		l.Printf("Strating the server on %s\n", bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error strating the server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the service
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until a signal is received
	sig := <-c
	log.Printf("Got signal: %s\n", sig)

	// gracefully shutdown the server
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	s.Shutdown(ctx)
}
