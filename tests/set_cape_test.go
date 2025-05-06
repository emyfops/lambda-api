package tests

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"github.com/Edouard127/lambda-api/api/models"
	"github.com/Edouard127/lambda-api/api/routes"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/go-redis/redismock/v9"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetCape(t *testing.T) {
	testCases := []struct {
		name        string
		token       string
		query       string
		preflight   func(mock redismock.ClientMock)
		expectError bool
	}{
		{
			name:  "Set user 1 cape",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7Im5hbWUiOiJ0ZXN0X3VzZXIxIiwiaWQiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDEiLCJkaXNjb3JkX2lkIjoiMDAwMDAwMDAwMDAwMDAwMDAxIiwidW5zYWZlIjpmYWxzZX0sImV4cCI6MjE0MTM4MDcwNCwiaWF0IjoxNzQxMjk0MzA0LCJuYmYiOjE3NDEyOTQzMDR9.VwctPUX2DzgsBnVxmlrtwexlPj3OQP4d0suGXttB6Mw",
			query: "?id=galaxy",
			preflight: func(mock redismock.ClientMock) {
				mock.ExpectSet("00000000-0000-0000-0000-000000000001", "galaxy", 0).SetVal("galaxy")
			},
			expectError: false,
		},
		{
			name:  "Set user 2 cape",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7Im5hbWUiOiJ0ZXN0X3VzZXIyIiwiaWQiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDIiLCJkaXNjb3JkX2lkIjoiMDAwMDAwMDAwMDAwMDAwMDAyIiwidW5zYWZlIjpmYWxzZX0sImV4cCI6MjE0MTM4MDcwNCwiaWF0IjoxNzQxMjk0MzA0LCJuYmYiOjE3NDEyOTQzMDR9.dyzi3eHRC2xB3nNmZKIuBjDvwBh4ADFP3F89Zvv2wFk",
			query: "?id=galaxy",
			preflight: func(mock redismock.ClientMock) {
				mock.ExpectSet("00000000-0000-0000-0000-000000000002", "galaxy", 0).SetVal("galaxy")
			},
			expectError: false,
		},
		{
			name:        "No algorithm",
			token:       "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJkYXRhIjp7Im5hbWUiOiJtYWxpY2lvdXNfdXNlciIsImlkIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwiZGlzY29yZF9pZCI6IjAwMDAwMDAwMDAwMDAwMDAwMCIsInVuc2FmZSI6ZmFsc2V9LCJleHAiOjIxNDEzODA3MDQsImlhdCI6MTc0MTI5NDMwNCwibmJmIjoxNzQxMjk0MzA0fQ.",
			expectError: true,
		},
		{
			name:        "Invalid query",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7Im5hbWUiOiJ0ZXN0X3VzZXIxIiwiaWQiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDEiLCJkaXNjb3JkX2lkIjoiMDAwMDAwMDAwMDAwMDAwMDAxIiwidW5zYWZlIjpmYWxzZX0sImV4cCI6MjE0MTM4MDcwNCwiaWF0IjoxNzQxMjk0MzA0LCJuYmYiOjE3NDEyOTQzMDR9.VwctPUX2DzgsBnVxmlrtwexlPj3OQP4d0suGXttB6Mw",
			query:       "?id=x",
			expectError: true,
		},
		{
			name:        "No login",
			expectError: true,
		},
	}

	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(true)

	app := fiber.New()
	internal.Set("logger", slog.Default())
	internal.Set("cache", db)
	app.Get("/", jwtware.New(jwtware.Config{
		SuccessHandler: func(ctx *fiber.Ctx) error {
			token := ctx.Locals("user").(*jwt.Token)

			parsed := token.Claims.(jwt.MapClaims)["data"]
			bytes, _ := json.Marshal(parsed)

			var player models.Player
			json.Unmarshal(bytes, &player)

			ctx.Locals("player", player)
			return ctx.Next()
		},
		// here we use hs256 because keys are deterministic
		SigningKey: jwtware.SigningKey{
			JWTAlg: "HS256",
			Key:    hmac.New(sha256.New, []byte{}).Sum(nil),
		},
	}), routes.SetCape)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.preflight != nil {
				tc.preflight(mock)
			}

			req := httptest.NewRequest(http.MethodGet, "/"+tc.query, nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)

			resp, err := app.Test(req)

			assert.Nil(t, err)

			if tc.expectError {
				assert.NotEqual(t, http.StatusOK, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}

			// All database calls must happen
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}

	// All database calls must happen
	assert.Nil(t, mock.ExpectationsWereMet())
}
