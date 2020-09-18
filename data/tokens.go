package data

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"
	"fmt"

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

// Token ...
type Token struct {
	ID        string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

////////// DATABASE //////////

var blackList map[string]bool

////////// FUNCTIONS //////////

// GenerateJWT ...
func GenerateJWT(id string) (string, error) {

	claims := &jwt.StandardClaims{
		IssuedAt: time.Now().Unix(),
		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		Id: id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// BlacklistToken ...
func BlacklistToken(tokenString string) {
	blackList[tokenString] = true
}

// ValidateToken ...
func ValidateToken(tokenString string) (*Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	t := &Token{
		ID: claims["jti"].(string),
		IssuedAt: time.Unix(int64(claims["iat"].(float64)), 0).UTC(),
		ExpiresAt: time.Unix(int64(claims["exp"].(float64)), 0).UTC(),
	}
	return t, nil
}