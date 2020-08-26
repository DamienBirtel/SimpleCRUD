package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// this is just for better readability
type username string

// User is the structure that holds information about a user
type User struct {
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

// Users is the map of users
type Users map[username]User

// ErrUserAlreadyCreated is used if we try to add an already created user to the db
var ErrUserAlreadyCreated = fmt.Errorf("User already created")

// ErrIncorrectLoginOrPassword expresses what it sounds like
var ErrIncorrectLoginOrPassword = fmt.Errorf("Incorrect Login or Password")

// usersList will act as our db of Users for now
var usersList = Users{
	"user": {
		Username:     "user",
		Password:     "password",
		RegisteredAt: time.Now().UTC(),
	},
}

// FromJSON reads data from an io.Reader and stores it in a User struct
func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// ToJSON encodes usersList to JSON and writes to w
func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// GetUsers returns the list of users for now
func GetUsers() Users {
	return usersList
}

// Register adds a User to the db, or returns an error if it already exists
func (u User) Register() error {

	// check to see if the user already exists, if so return an error
	if _, ok := usersList[username(u.Username)]; ok {
		return ErrUserAlreadyCreated
	}

	// add the time when we create the user and add it to the db
	u.RegisteredAt = time.Now().UTC()
	usersList[username(u.Username)] = u
	return nil
}

// Login tries to login a user and return an error if it fails
func (u User) Login() error {

	existingUser, ok := usersList[username(u.Username)]
	if !ok {
		return ErrIncorrectLoginOrPassword
	}

	if u.Password != existingUser.Password {
		return ErrIncorrectLoginOrPassword
	}

	return nil
}

func (u User) Delete() {
	delete(usersList, username(u.Username))
}
