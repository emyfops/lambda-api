package random

import (
	"testing"
)

func TestRandString(t *testing.T) {
	t.Log(RandString(69))
	t.Log(RandString(69))
	t.Log(RandString(69))
}

func TestRandInt32(t *testing.T) {
	t.Log(RandInt32())
	t.Log(RandInt32())
	t.Log(RandInt32())
}

func TestRandInt64(t *testing.T) {
	t.Log(RandInt64())
	t.Log(RandInt64())
	t.Log(RandInt64())
}

func TestRandFloat(t *testing.T) {
	t.Log(RandFloat())
	t.Log(RandFloat())
	t.Log(RandFloat())
}

func TestRandDouble(t *testing.T) {
	t.Log(RandDouble())
	t.Log(RandDouble())
	t.Log(RandDouble())
}
