package tests

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Edouard127/lambda-api/api/models"
	"github.com/Edouard127/lambda-api/api/routes"
	"github.com/Edouard127/lambda-api/internal"
	"github.com/go-redis/redismock/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	request := models.CapeLookup{Players: make([]uuid.UUID, 0)}

	for i := 0; i < 1; i++ {
		id, _ := uuid.NewRandom()
		request.Players = append(request.Players, id)
	}

	b, _ := json.Marshal(request)

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/capes", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			http.DefaultClient.Do(req)
		}()
	}

	wg.Wait()
}

func TestGetCape(t *testing.T) {
	testCases := []struct {
		name        string
		query       string
		preflight   func(mock redismock.ClientMock)
		expectError bool
	}{
		{
			name:  "Player id",
			query: "?id=00000000-0000-0000-0000-000000000000",
			preflight: func(mock redismock.ClientMock) {
				mock.ExpectGet("00000000-0000-0000-0000-000000000000").SetVal("1")
			},
			expectError: false,
		},
		{
			name:  "Additional player ids",
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
	internal.Set("key", &rsa.PrivateKey{})
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
