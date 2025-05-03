package tests

import (
	"github.com/Edouard127/lambda-api/api/routes"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/go-redis/redismock/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCape(t *testing.T) {
	testCases := []struct {
		name        string
		query       string
		preflight   func(mock redismock.ClientMock)
		expectError bool
	}{
		{
			name:  "Locals player id",
			query: "?id=00000000-0000-0000-0000-000000000000",
			preflight: func(mock redismock.ClientMock) {
				mock.ExpectGet("00000000-0000-0000-0000-000000000000").SetVal("1")
			},
			expectError: false,
		},
		{
			name:  "Locals additional player ids",
			query: "?id=00000000-0000-0000-0000-000000000000&id=ab24f5d6-dcf1-45e4-897e-b50a7c5e7422",
			preflight: func(mock redismock.ClientMock) {
				mock.ExpectGet("00000000-0000-0000-0000-000000000000").SetVal("1")
			},
			expectError: false,
		},
		{
			name:        "Without player id",
			query:       "",
			expectError: true,
		},
	}

	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(true)

	app := fiber.New()
	internal.Set("logger", slog.Default())
	internal.Set("cache", db)
	app.Get("/", routes.GetCape)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.preflight != nil {
				tc.preflight(mock)
			}

			req := httptest.NewRequest(http.MethodGet, "/"+tc.query, nil)

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
}
