package handlers

import (
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
