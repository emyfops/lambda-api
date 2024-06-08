package memory

import "testing"

func TestNewCache(t *testing.T) {
	c := NewCache[string, int]()
	if c == nil {
		t.Error("expected cache to be initialized")
	}

	if c.items == nil {
		t.Error("expected cache items to be initialized")
	}

	for i := 0; i < 9000; i++ {
		c.Set(string(rune(i)), i, -1)
	}

	for i := 0; i < 9000; i++ {
		_, exists := c.Get(string(rune(i)))
		if !exists {
			t.Errorf("expected key %d to exist", i)
		}
	}

	for i := 0; i < 9000; i++ {
		c.Delete(string(rune(i)))
		_, exists := c.Get(string(rune(i)))
		if exists {
			t.Errorf("expected key %d to not exist", i)
		}
	}

	if len(c.items) != 0 {
		t.Error("expected cache to be empty")
	}

	t.Log("TestNewCache passed")
}
