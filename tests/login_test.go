package tests

import (
	"bytes"
	"encoding/json"
	"github.com/Edouard127/lambda-api/api/middlewares"
	"github.com/Edouard127/lambda-api/api/models"
	"github.com/Edouard127/lambda-api/api/routes"
	"github.com/gofiber/fiber/v2"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Parallel()
	flag.Set("online", "false")

	testCases := []struct {
		name        string
		body        models.Authentication
		expectError bool
	}{
		{
			name: "Locals invalid name",
			body: models.Authentication{
				Username: "-",
				Hash:     "",
			},
			expectError: true,
		},
		{
			name: "Locals invalid name",
			body: models.Authentication{
				Username: "ｕｓｅｒｎａｍｅ２",
				Hash:     "",
			},
			expectError: true,
		},
		{
			name: "Locals invalid name",
			body: models.Authentication{
				Username: "u̸̥͍͂̐͐ͅs̸̨̝̈ȅ̶̫̚r̷̠̆̽͘͝n̶͚̓̓̈́͝a̸͍̱͇̒͂m̵̞̦̰̽̇̔͝e̷̓͜3̵̛͖̱̗̗̌",
				Hash:     "",
			},
			expectError: true,
		},
	}

	app := fiber.New()
	app.Use(middlewares.Locals("logger", slog.Default()))
	app.Get("/", routes.Login)

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
