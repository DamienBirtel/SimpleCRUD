package data

import (
	"encoding/json"
	"io"
	"time"
)

type username string

// User is the structure that holds information about a user
type User struct {
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

// Users is the map of users
type Users map[username]User

var usersList = Users{
	"user": {
		Username:     "user",
		Password:     "password",
		RegisteredAt: time.Now().UTC(),
	},
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
