package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DamienBirtel/SimpleCRUD/data"
)

// Handler handles the requests and logs info to a logger
type Handler struct {
	l *log.Logger
}

// NewHandler returns a new Handler{}
func NewHandler(l *log.Logger) *Handler {
	return &Handler{l}
}

// SignUp adds a new data.User{} to the db
func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(KeyUser{}).(*data.User)

	h.l.Println("[INFO] Adding User: ", u.Username)

	err := data.AddUser(u)
	if err != nil {
		h.l.Println("[ERROR] Can't add User: ", u.Username, err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	resp := fmt.Sprintf("%s successfully signed up !", u.Username)
	w.Write([]byte(resp))
}

// GetUsers returns the users as a JSON
func (h Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := data.GetUsers()
	users.ToJSON(w)
}
