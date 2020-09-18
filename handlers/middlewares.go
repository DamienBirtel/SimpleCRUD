package handlers

import (
	"context"
	"net/http"
	"strings"
	"github.com/DamienBirtel/SimpleCRUD/data"
)

// KeyUser is used to pass data.User{} info to the request
type KeyUser struct{}

// KeyToken is used to pass data.Token{} info to the request
type KeyToken struct{}

// MiddlewareValidateUser validates that a data.User{} is sent through the request
func (h Handler) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		u := &data.User{}

		err := u.FromJSON(r.Body)
		if err != nil {
			h.l.Println("[ERROR] deserializing user ", u, err)
			http.Error(w, "Error reading User info", http.StatusBadRequest)
			return
		}

		err = u.Validate()
		if err != nil {
			h.l.Println("[ERROR] validating user ", u, err)
			http.Error(w, "Invalid username or password", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, u)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// MiddlewareValidateToken ...
func (h Handler) MiddlewareValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authString := r.Header.Get("Authorization")
		if authString == "" {
			h.l.Println("[ERROR] no token")
			http.Error(w, "You need to log in to access this ressource", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authString, " ")[1]

		token, err := data.ValidateToken(tokenString)
		if err != nil {
			h.l.Println("[ERROR] validating token", tokenString)
			http.Error(w, "You need to log in to access this ressource", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), KeyToken{}, token)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}