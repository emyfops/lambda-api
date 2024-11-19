package auth

import (
	"testing"
)

func TestNewJwt(t *testing.T) {
	type test struct {
		Info string `json:"info"`
	}

	signed, err := NewJwt(&test{Info: "aaa"})
	if err != nil {
		t.Error(err)
	}

	token, err := ParseString(signed)
	if err != nil {
		t.Error(err)
	}

	var result test
	ParseStruct(token, &result)
	if result.Info != "aaa" {
		t.Errorf("expected 'aaa', got %s", result.Info)
		return
	}

	t.Log("TestNewJwt passed")
}
