package util

import (
	"fmt"
	"testing"
)

func TestNewJwt(t *testing.T) {
	t.Run("TestNewJwt", func(t *testing.T) {
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
		JwtToStruct(token, &result)
		fmt.Println(result)
	})
}
