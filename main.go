package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {

}

func createRouter() *mux.Router {

	serveMux := mux.NewRouter()
	serveMux.HandleFunc("/", handler)
	return serveMux
}

func createServer() *http.Server {

	serveMux := createRouter()

	// 	TODO: MAKE THIS CONFIGURABLE !!!
	srv := &http.Server{
		Addr:         "localhost:9090",
		Handler:      serveMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}

	return srv
}

func main() {

	// We first create our server
	srv := createServer()

	// We plan the graceful shutdown of the server by catching interrupt signal
	idleConnsClosed := make(chan struct{})
	go func() {
		defer close(idleConnsClosed)
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, os.Kill)
		<-sigint

		// We received an interrupt signal, now we shut down.
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Printf("HHTP server Shutdown error: %v", err)
		}
	}()

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		// Error starting or closing listener
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
}
