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

func createRouter() *mux.Router {

	l := log.New(os.Stdout, "SimpleCRUD", log.LstdFlags)
	h := handlers.NewHandler(l)

	router := mux.NewRouter()

	listR := router.Methods(http.MethodGet).Subrouter()
	listR.HandleFunc("/list", h.List)

	getR := router.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/", h.Get)
	getR.HandleFunc("/logout", h.Logout)
	getR.Use(h.MiddlewareValidateToken)

	postR := router.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/sign_up", h.SignUp)
	postR.HandleFunc("/sign_in", h.SignIn)
	postR.Use(h.MiddlewareValidateUser)

	putR := router.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/update_username", h.UpdateUsername)
	putR.HandleFunc("/update_password", h.UpdatePassword)
	putR.Use(h.MiddlewareValidateToken)
	putR.Use(h.MiddlewareValidateUser)

	delR := router.Methods(http.MethodDelete).Subrouter()
	delR.HandleFunc("/delete", h.Delete)
	delR.Use(h.MiddlewareValidateToken)

	return router
}

func createServer() *http.Server {

	router := createRouter()

	srv := &http.Server{
		Addr:         "localhost:9090",
		Handler:      router,
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
			log.Printf("HTTP server Shutdown error: %v", err)
		}
	}()

	log.Println("Starting server...")
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		// Error starting or closing listener
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
}
