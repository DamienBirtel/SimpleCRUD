package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/DamienBirtel/SimpleCRUD/data"
)

// Get is the basic handler, and should for now give us a list of all registered users
func Get(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("id")
	rw.Write([]byte(fmt.Sprintf("Hello %s !", user)))
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

	// we get a data.User object through context
	u := r.Context().Value(KeyUser{}).(data.User)

	// we check if the username and the password match
	err := u.Login()
	if err != nil {
		http.Error(rw, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	// if successful return a jwt token containing the user ID
	token, err := data.NewJWT(u.Username)
	if err != nil {
		http.Error(rw, "Error creating new token", http.StatusInternalServerError)
		return
	}
	rw.Write([]byte(token))
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

// Logout logs users out
func Logout(rw http.ResponseWriter, r *http.Request) {

}

// Update ...
func Update(rw http.ResponseWriter, r *http.Request) {

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

// MiddlewareAuthenticate checks for a valid jwt token in the request header
func MiddlewareAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(tokenString) != 2 {
			rw.Write([]byte("You are not logged in\n"))
			return
		}

		token, err := data.ValidateJWT(tokenString[1])
		if err != nil {
			http.Error(rw, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "id", token.Id)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
