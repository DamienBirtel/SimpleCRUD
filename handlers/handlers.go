package handlers

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

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

// GetUsers returns the users as a JSON
func (h Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := data.GetUsers()
	users.ToJSON(w)
}

// SignUp adds a new data.User{} to the db
func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(KeyUser{}).(*data.User)

	err := data.AddUser(u)
	if err != nil {
		h.l.Println("[ERROR] Can't sign up User: ", u.Username, err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	h.l.Println("[INFO] Signing up User: ", u.Username)

	resp := fmt.Sprintf("%s successfully signed up !", u.Username)
	w.Write([]byte(resp))
}

// SignIn returns a jwt
func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {

	signInUser := r.Context().Value(KeyUser{}).(*data.User)

	dbUser, err := data.GetUser(signInUser.Username)
	if err != nil {
		h.l.Println("[WARNING] Can't sign in User: ", signInUser.Username, err)
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(signInUser.Password))
	if err != nil {
		h.l.Println("[WARNING] Can't sign in User: ", signInUser.Username, err)
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	h.l.Println("[INFO] Signing in User: ", signInUser.Username)
	w.Write([]byte("You are now logged in !"))
}
