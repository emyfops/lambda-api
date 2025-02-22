package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var secretKey []byte

func init() {
	mac := hmac.New(sha256.New, []byte(os.Getenv("SECRET_KEY")))
	secretKey = mac.Sum(nil)
}

// NewJwt generates a JWT and signs it with a certificate
func NewJwt(claims any) (signed string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf":  time.Now().Unix(),
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"data": claims,
	})

	signed, err = token.SignedString(secretKey)
	if err != nil {
		err = errors.New("failed to sign the token: " + err.Error())
		return
	}

	return
}

// ParseJwt parses a signed JWT token from a given string
func ParseJwt(signed string) (token *jwt.Token, err error) {
	token, err = jwt.Parse(signed, func(token *jwt.Token) (any, error) {
		// Avoid the none signing method attack
		_, valid := token.Method.(*jwt.SigningMethodHMAC)
		if !valid {
			return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
		}

		return secretKey, nil
	})

	if err != nil {
		return token, errors.New("failed to parse the token: " + err.Error())
	}

	return token, nil
}

// ParseStructJwt takes a given JWT token and parses the data field into the provided struct pointer
func ParseStructJwt[T any](token *jwt.Token, result *T) error {
	parsed, ok := token.Claims.(jwt.MapClaims)["data"]
	if !ok {
		return errors.New("data field not found in JWT claims")
	}

	bytes, err := json.Marshal(parsed)
	if err != nil {
		return errors.New("error marshalling data field:" + err.Error())
	}

	err = json.Unmarshal(bytes, result)
	if err != nil {
		return errors.New("error parsing data into provided struct: " + err.Error())
	}

	return nil
}
