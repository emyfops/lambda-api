package tests

import (
	"github.com/Edouard127/lambda-api/internal/app/jwt"
	"testing"
)

func TestNewJwt(t *testing.T) {
	type test struct {
		Info string `json:"info"`
	}

	signed, err := jwt.New(&test{Info: "aaa"})
	if err != nil {
		t.Error(err)
	}

	token, err := jwt.ParseString(signed)
	if err != nil {
		t.Error(err)
	}

	var result test
	jwt.ParseStruct(token, &result)
	if result.Info != "aaa" {
		t.Errorf("expected 'aaa', got %s", result.Info)
		return
	}

	t.Log("TestNewJwt passed")
}
