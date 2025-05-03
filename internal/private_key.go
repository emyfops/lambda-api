package internal

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var PrivateKey = loadPrivateKeyFromFile("key/id_rsa")

func loadPrivateKeyFromFile(path string) *rsa.PrivateKey {
	f, err := os.ReadFile(path)
	if err != nil {
		return genKey()
	}

	k, err := jwt.ParseRSAPrivateKeyFromPEM(f)
	if err != nil {
		return genKey()
	}

	return k
}

func genKey() *rsa.PrivateKey {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	return key
}
