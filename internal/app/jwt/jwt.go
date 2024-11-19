package jwt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var privKey *rsa.PrivateKey
var pubKey *rsa.PublicKey
var pubKeyBytes []byte
var pubKeyBlock *pem.Block
var pubKeyBuffer bytes.Buffer

func init() {
	var err error

	privKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("failed to generate private key: " + err.Error())
	}

	pubKey = &privKey.PublicKey

	pubKeyBytes, err = x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		panic("failed to marshal public key: " + err.Error())
	}

	pubKeyBlock = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	if err = pem.Encode(&pubKeyBuffer, pubKeyBlock); err != nil {
		panic("failed to encode public key: " + err.Error())
	}
}

// New generates a JWT and signs it with a certificate
func New(claims any) (signed string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"nbf":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"aud":  "lotsOfUsers",
		"data": claims,
	})

	signed, err = token.SignedString(privKey)
	if err != nil {
		err = errors.New("failed to sign token: " + err.Error())
		return
	}

	return
}

// ParseString parses a signed JWT token from a given string
func ParseString(signed string) (token *jwt.Token, err error) {
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return token, errors.New("failed to parse public key: " + err.Error())
	}

	token, err = jwt.Parse(signed, func(token *jwt.Token) (any, error) { return pubKey, nil })
	if err != nil {
		return token, errors.New("failed to parse token: " + err.Error())
	}

	return token, nil
}

// ParseStruct takes a given JWT token and parses the data field into the provided struct pointer
func ParseStruct[T any](token *jwt.Token, result *T) error {
	parsed, ok := token.Claims.(jwt.MapClaims)["data"]
	if !ok {
		return errors.New("data field not found in JWT claims")
	}

	dataBytes, err := json.Marshal(parsed)
	if err != nil {
		return errors.New("error marshalling data field:" + err.Error())
	}

	if err = json.Unmarshal(dataBytes, result); err != nil {
		return errors.New("error unmarshalling data field into provided struct: " + err.Error())
	}

	return nil
}
