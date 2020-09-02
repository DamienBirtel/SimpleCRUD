package handlers

import (
	"context"
	"net/http"

	"github.com/DamienBirtel/SimpleCRUD/data"
)

// KeyUser is used to pass data.User{} info to the request
type KeyUser struct{}

// MiddleWareValidateUser validates that a data.User{} is sent through the request
func (h Handler) MiddleWareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		u := &data.User{}

		err := u.FromJSON(r.Body)
		if err != nil {
			h.l.Println("[ERROR] deserializing user ", u, err)
			http.Error(w, "Error reading User info", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, u)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
