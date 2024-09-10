package auth

import (
	"testing"
)

func TestNewJwt(t *testing.T) {
	type test struct {
		Info string `json:"info"`
	}

	signed, err := CreateJwtToken(&test{Info: "aaa"})
	if err != nil {
		t.Error(err)
	}

	token, err := ParseJwtToken(signed)
	if err != nil {
		t.Error(err)
	}

	var result test
	ParseToStruct(token, &result)
	if result.Info != "aaa" {
		t.Errorf("expected 'aaa', got %s", result.Info)
		return
	}

	t.Log("TestNewJwt passed")
}
