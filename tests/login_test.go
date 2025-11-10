package tests

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Edouard127/lambda-api/api/models"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/gofiber/fiber/v2"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	flag.Set("online", "false")

	testCases := []struct {
		name        string
		body        models.Authentication
		expectError bool
	}{
		{
			name: "Invalid name",
			body: models.Authentication{
				Username: "-",
				Hash:     "",
			},
			expectError: true,
		},
		{
			name: "Invalid name",
			body: models.Authentication{
				Username: "ｕｓｅｒｎａｍｅ２",
				Hash:     "",
			},
			expectError: true,
		},
		{
			name: "Invalid name",
			body: models.Authentication{
				Username: "u̸̥͍͂̐͐ͅs̸̨̝̈ȅ̶̫̚r̷̠̆̽͘͝n̶͚̓̓̈́͝a̸͍̱͇̒͂m̵̞̦̰̽̇̔͝e̷̓͜3̵̛͖̱̗̗̌",
				Hash:     "",
			},
			expectError: true,
		},
	}

	app := fiber.New()
	internal.Set("logger", slog.Default())
	internal.Set("key", &rsa.PrivateKey{})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			assert.Nil(t, err)

			if tc.expectError {
				assert.NotEqual(t, http.StatusOK, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}
		})
	}
}
