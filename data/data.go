package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User is the structure that holds information about a user
type User struct {
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

// Users is the map of users
type Users map[string]User

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

// blacklist lists all expired tokens
var blacklist map[string]int

var secretKey = []byte("super secret key")

// NewJWT returns a new JWT
func NewJWT(username string) (string, error) {

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		Id:        username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type LittleToken struct {
	Id string
}

// ValidateJWT ...
func ValidateJWT(tokenString string) (*LittleToken, error) {

	// we first check to see if the string isn't blacklisted
	if _, ok := blacklist[tokenString]; ok {
		return nil, fmt.Errorf("Blacklisted token")
	}

	// we parse the token string to get a token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// we make sure the signing method is the right one
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	lt := &LittleToken{
		Id: claims["jti"].(string),
	}
	// TODO implement expires at check //

	return lt, nil
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
	if _, ok := usersList[u.Username]; ok {
		return ErrUserAlreadyCreated
	}

	// add the time when we create the user and add it to the db
	u.RegisteredAt = time.Now().UTC()
	usersList[u.Username] = u
	return nil
}

// Login tries to login a user and return an error if it fails
func (u User) Login() error {

	existingUser, ok := usersList[u.Username]
	if !ok {
		return ErrIncorrectLoginOrPassword
	}

	if u.Password != existingUser.Password {
		return ErrIncorrectLoginOrPassword
	}

	return nil
}

// Delete deletes the user
func (u User) Delete() {
	delete(usersList, u.Username)
}
