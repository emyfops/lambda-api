package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
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
		panic(fmt.Sprintf("failed to generate key: %s\n", err))
	}

	pubKey = &privKey.PublicKey

	pubKeyBytes, err = x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal public key: %s\n", err))
	}

	pubKeyBlock = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	if err := pem.Encode(&pubKeyBuffer, pubKeyBlock); err != nil {
		panic(fmt.Sprintf("failed to encode public key: %s\n", err))
	}
}

// CreateJwtToken generates a JWT (and signs it with a generated certificate)
func CreateJwtToken(claims any) (signed string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"nbf":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"aud":  "lotsOfUsers",
		"data": claims,
	})

	signed, err = token.SignedString(privKey)
	if err != nil {
		err = fmt.Errorf("failed to sign token: %s\n", err)
		return
	}

	return
}

// ParseJwtToken parses a JWT token
func ParseJwtToken(signed string) (token *jwt.Token, err error) {
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return token, fmt.Errorf("failed to parse public key: %s\n", err)
	}

	token, err = jwt.Parse(signed, func(token *jwt.Token) (any, error) {
		return pubKey, nil
	})
	if err != nil {
		return token, fmt.Errorf("failed to parse token: %s\n", err)
	}

	return token, nil
}

func ParseToStruct[T any](token *jwt.Token, result *T) error {
	parsed, ok := token.Claims.(jwt.MapClaims)["data"]
	if !ok {
		return errors.New("data field not found in JWT claims")
	}

	dataBytes, err := json.Marshal(parsed)
	if err != nil {
		return fmt.Errorf("error marshalling data field: %v", err)
	}

	if err := json.Unmarshal(dataBytes, result); err != nil {
		return fmt.Errorf("error unmarshalling data field into provided struct: %v", err)
	}

	return nil
}
