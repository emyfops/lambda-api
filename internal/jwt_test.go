package internal

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestParseStructJwt(t *testing.T) {
	now := time.Now()
	type args[T any] struct {
		token  *jwt.Token
		result *T
	}
	type testCase[T any] struct {
		name    string
		args    args[T]
		wantErr bool
	}
	tests := []testCase[struct {
		n0  int
		n1  int
		now time.Time
	}]{
		{
			name: "ok",
			args: args[struct {
				n0  int
				n1  int
				now time.Time
			}]{
				token: jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"data": jwt.MapClaims{
						"n0":  0,
						"n1":  1,
						"now": now,
					},
				}),
				result: &struct {
					n0  int
					n1  int
					now time.Time
				}{
					n0:  0,
					n1:  1,
					now: now,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseStructJwt(tt.args.token, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("ParseStructJwt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
