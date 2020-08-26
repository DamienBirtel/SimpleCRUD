package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DamienBirtel/SimpleCRUD/data"
)

// Get is the basic handler, and should for now give us a list of all registered users
func Get(rw http.ResponseWriter, r *http.Request) {
	users := data.GetUsers()
	err := users.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

// Register gets the user info from the request body and adds it to the db
func Register(rw http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(KeyUser{}).(data.User)

	err := u.Register()
	if err == data.ErrUserAlreadyCreated {
		http.Error(rw, "User already exists", http.StatusBadRequest)
		return
	}

	// let them know everything went fine
	rw.Write([]byte("User successfully created"))
}

// Login checks if the user exists and if the password matches and... TODO implement tokenization
func Login(rw http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(KeyUser{}).(data.User)

	err := u.Login()
	if err != nil {
		http.Error(rw, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	rw.Write([]byte(fmt.Sprintf("%s successfully Logged in\n", u.Username)))
}

// Delete deletes a user
func Delete(rw http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(KeyUser{}).(data.User)

	err := u.Login()
	if err != nil {
		http.Error(rw, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	u.Delete()

	rw.Write([]byte(fmt.Sprintf("User: %s, successfully deleted\n", u.Username)))
}

// KeyUser yay
type KeyUser struct{}

// MiddlewareValidateUserInfo checks if the request is formatted correctly
// before sending it to another handler
func MiddlewareValidateUserInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		u := data.User{}

		// we try to convert what's in the request body to a user object and return if it fails
		err := u.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
			return
		}

		// if not we add it to the request through context
		ctx := context.WithValue(r.Context(), KeyUser{}, u)
		r = r.WithContext(ctx)

		// we call the next handler in the chain
		next.ServeHTTP(rw, r)
	})
}
