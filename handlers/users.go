package handlers

import (
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

	// create an empty User struct
	u := data.User{}

	// try to fill it with info from the request body
	err := u.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	// try to add the user to the database
	err = u.Register()
	if err == data.ErrUserAlreadyCreated {
		http.Error(rw, "User already exists", http.StatusBadRequest)
		return
	}

	// let them know everything went fine
	rw.Write([]byte("User successfully created"))
}

// Login checks if the user exists and if the password matches and... TODO implement tokenization
func Login(rw http.ResponseWriter, r *http.Request) {

	u := data.User{}

	err := u.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	err = u.Login()
	if err != nil {
		http.Error(rw, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	rw.Write([]byte(fmt.Sprintf("%s successfully Logged in\n", u.Username)))
}

// Delete deletes a user
func Delete(rw http.ResponseWriter, r *http.Request) {

	u := data.User{}

	err := u.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	err = u.Login()
	if err != nil {
		http.Error(rw, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	u.Delete()

	rw.Write([]byte(fmt.Sprintf("User: %s, successfully deleted\n", u.Username)))
}
