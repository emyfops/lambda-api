package tests

import (
	"github.com/Edouard127/lambda-api/internal"
	"testing"
)

func TestNewJwt(t *testing.T) {
	type test struct {
		Info string `json:"info"`
	}

	signed, err := internal.NewJwt(&test{Info: "aaa"})
	if err != nil {
		t.Error(err)
	}

	token, err := internal.ParseJwt(signed)
	if err != nil {
		t.Error(err)
	}

	var result test
	internal.ParseStructJwt(token, &result)
	if result.Info != "aaa" {
		t.Errorf("expected 'aaa', got %s", result.Info)
		return
	}

	t.Log("TestNewJwt passed")
}
