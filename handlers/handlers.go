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

// Get shows some info about the logged in user
func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	
	t := r.Context().Value(KeyToken{}).(*data.Token)

	// now := time.Now()

	// timeSinceLogIn := int(now.Sub(t.IssuedAt).Minutes())

	// expelTime := int(t.ExpiresAt.Sub(now).Minutes())

	fmt.Fprintf(w, "Hello %s !", t.ID)// ! You have been logged in since %d minutes and will be automatically logged out in %d minutes", t.ID, timeSinceLogIn, expelTime)
}

// SignUp adds a new data.User{} to the db
func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {

	u := r.Context().Value(KeyUser{}).(*data.User)

	err := data.AddUser(u)
	if err != nil {
		h.l.Println("[ERROR] Can't sign up User:", u.Username, err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	h.l.Println("[INFO] Signing up User:", u.Username)

	fmt.Fprintln(w, "Congratulation on signing up, you can now sign in !")
}

// SignIn returns a jwt
func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {

	signInUser := r.Context().Value(KeyUser{}).(*data.User)

	dbUser, err := data.GetUser(signInUser.Username)
	if err != nil {
		h.l.Println("[WARNING] Can't sign in User:", signInUser.Username, err)
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(signInUser.Password))
	if err != nil {
		h.l.Println("[WARNING] Can't sign in User:", signInUser.Username, err)
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	tokenString, err := data.GenerateJWT(signInUser.Username)
	if err != nil {
		h.l.Println("[ERROR] Can't sign token:", err)
		http.Error(w, "Error while signing token", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, tokenString)
}

// Logout logs the user out, which means it blacklists the affiliated token
func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	
	t := r.Context().Value(KeyToken{}).(*data.Token)

	data.BlacklistToken(t.TokenString)
}

// UpdateUsername ...
func (h Handler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	
	u := r.Context().Value(KeyUser{}).(*data.User)
	t := r.Context().Value(KeyToken{}).(*data.Token)

	err := data.UpdateUserName(t.ID, u.Username)
	if err != nil {
		if err == data.ErrUsernameAlreadyTaken {
			http.Error(w, "Username already taken", http.StatusNotAcceptable)
			return
		}
		h.l.Println("[ERROR] ", err)
		return
	}

	tokenString, err := data.GenerateJWT(u.Username)
	if err != nil {
		h.l.Println("[ERROR] Can't sign token:", err)
		http.Error(w, "Error while signing token", http.StatusInternalServerError)
		return
	}

	data.BlacklistToken(t.TokenString)
	fmt.Fprintf(w, "Your username has been updated to %s, here is your new token: %s", u.Username, tokenString)
}

// UpdatePassword ...
func (h Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	
	u := r.Context().Value(KeyUser{}).(*data.User)
	t := r.Context().Value(KeyToken{}).(*data.Token)

	err := data.UpdateUserPassword(t.ID, u.Password)
	if err != nil {
		h.l.Println("[ERROR] ", err)
		http.Error(w, "Forbidden", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Your password has been updated")
}

// Delete ...
func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {

	t := r.Context().Value(KeyToken{}).(*data.Token)
	
	data.BlacklistToken(t.TokenString)
	data.DeleteUser(t.ID)
	fmt.Fprintf(w, "Your account has been deleted")
}

// List ...
func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	users := data.GetUsers()
	users.ToJSON(w)
}