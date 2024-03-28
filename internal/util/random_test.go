package util

import (
	"fmt"
	"testing"
)

func TestRandStringBytesMaskSrcUnsafe(t *testing.T) {
	fmt.Println(RandStringBytesMaskSrcUnsafe(10))
	fmt.Println(RandStringBytesMaskSrcUnsafe(16))
	fmt.Println(RandStringBytesMaskSrcUnsafe(77))
}
