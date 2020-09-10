package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

////////// STRUCTURE //////////

// User holds information about a User
type User struct {
	Username     string    `json:"username" validate:"required"`
	Password     string    `json:"password" validate:"required"`
	RegisteredAt time.Time `json:"registered_at"`
}

/// METHODS ///

// FromJSON decodes json from an io.Reader to the User{}
func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// ToJSON encodes json from the User{} to an io.Writer
func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// Validate returns an error if the struct doesn't match the validate fields
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

////////// DATABASE //////////

// Users is a custom type created to handle th ToJSON method
type Users map[string]*User

var usersList = Users{
	"John": {
		Username:     "John",
		Password:     "Password",
		RegisteredAt: time.Now().UTC(),
	},
}

/// METHODS ///

// ToJSON encodes the Users to the io.Writer
func (us Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(us)
}

/// ERRORS ///

// ErrUserNotFound is used when trying to get an unexisting user from the db
var ErrUserNotFound = fmt.Errorf("User not found")

// ErrUsernameAlreadyTaken is used when trying to create an already existing user
var ErrUsernameAlreadyTaken = fmt.Errorf("Username already taken")

/// FUNCTIONS ///

// GetUsers returns the usersList
func GetUsers() Users {
	return usersList
}

// GetUser returns a User{} from the db
func GetUser(username string) (*User, error) {

	u, ok := usersList[username]
	if !ok {
		return nil, ErrUserNotFound
	}

	return u, nil
}

// AddUser adds a User{} to the db
func AddUser(u *User) error {

	_, ok := usersList[u.Username]
	if ok {
		return ErrUsernameAlreadyTaken
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
	if err != nil {
		return err
	}

	u.Password = string(pwd)
	u.RegisteredAt = time.Now().UTC()
	usersList[u.Username] = u
	return nil

}

// DeleteUser deletes a User{} from the db
func DeleteUser(username string) {
	delete(usersList, username)
}

// UpdateUserPassword updates the password for a User{}
func UpdateUserPassword(username string, newPassword string) error {

	u, ok := usersList[username]
	if !ok {
		return ErrUserNotFound
	}

	u.Password = newPassword
	usersList[username] = u
	return nil
}

// UpdateUserName updates the username of a User
func UpdateUserName(oldUsername string, newUsername string) error {

	u, ok := usersList[oldUsername]
	if !ok {
		return ErrUserNotFound
	}

	u.Username = newUsername
	err := AddUser(u)
	if err != nil {
		return ErrUsernameAlreadyTaken
	}

	DeleteUser(oldUsername)
	return nil
}
