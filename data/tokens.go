package data

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

////////// KEYS //////////

const (
	privateKeyPath = "./keys/key.rsa"
	publicKeyPath  = "./keys/key.rsa.pub"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func init() {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
	}
}

// RefreshToken ...
type RefreshToken struct {
	ID        string
	ExpiresAt time.Time
}

////////// DATABASE //////////

var blackList = map[string]RefreshToken{}

////////// FUNCTIONS //////////

// GenerateJWT ...
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	claims["testlol"] = "John"

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
