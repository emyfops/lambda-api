package tests

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Edouard127/lambda-api/api/routes"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/go-redis/redismock/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetCapes(t *testing.T) {
	testCases := []struct {
		name        string
		body        map[string]interface{}
		preflight   func(mock redismock.ClientMock)
		expectError bool
	}{
		{
			name: "Get player ids",
			body: map[string]interface{}{
				"players": []string{
					"00000000-0000-0000-0000-000000000000",
					"00000000-0000-0000-0000-000000000001",
					"00000000-0000-0000-0000-000000000002",
					"00000000-0000-0000-0000-000000000003",
					"00000000-0000-0000-0000-000000000004",
				},
			},
			preflight: func(mock redismock.ClientMock) {
				mock.ExpectMGet(
					"00000000-0000-0000-0000-000000000000",
					"00000000-0000-0000-0000-000000000001",
					"00000000-0000-0000-0000-000000000002",
					"00000000-0000-0000-0000-000000000003",
					"00000000-0000-0000-0000-000000000004",
				).SetVal([]any{"0", "1", "2", "3", "4"})
			},
			expectError: false,
		},
		{
			name: "Get player ids with non present ids",
			body: map[string]interface{}{
				"players": []string{
					"00000000-0000-0000-0000-000000000000",
					"00000000-0000-0000-0000-000000000001",
					"00000000-0000-0000-0000-000000000002",
					"00000000-0000-0000-0000-000000000003",
					"00000000-0000-0000-0000-000000000004",
					"00000001-0000-0000-0000-000000000000",
				},
			},
			preflight: func(mock redismock.ClientMock) {
				mock.ExpectMGet(
					"00000000-0000-0000-0000-000000000000",
					"00000000-0000-0000-0000-000000000001",
					"00000000-0000-0000-0000-000000000002",
					"00000000-0000-0000-0000-000000000003",
					"00000000-0000-0000-0000-000000000004",
					"00000001-0000-0000-0000-000000000000",
				).SetVal([]any{"0", "1", "2", "3", "4", nil})
			},
			expectError: false,
		},
		{
			name: "Wrong player ids",
			body: map[string]interface{}{
				"players": []string{
					"00000000-0000-0000-0000-000000000000",
					"x",
				},
			},
			expectError: true,
		},
		{
			name:        "No player ids",
			body:        make(map[string]interface{}),
			expectError: false,
		},
	}

	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(true)

	app := fiber.New()
	internal.Set("logger", slog.Default())
	internal.Set("cache", db)
	internal.Set("key", &rsa.PrivateKey{})

	app.Get("/", routes.GetCapes)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.preflight != nil {
				tc.preflight(mock)
			}

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

	// All database calls must happen
	assert.Nil(t, mock.ExpectationsWereMet())
}
