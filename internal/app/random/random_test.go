package random

import (
	"math/rand"
	"testing"
)

func init() {
	src = rand.NewSource(0)
}

func TestRandString(t *testing.T) {
	const (
		first  = "lwDqHvrrrpHixIlOLIHafuuyunSEcEPIVngpTcYkBVsROQvIpxLpBSHZBbblhaJlRczqb"
		second = "tdiZlIZDFmNRbELhbZmnJfzvVIETOlZmPXcXkYZfMATHBoiJJVSCzetWxHmTpNvRTdqjo"
		third  = "VnXybyfsIHoPSZEGRhiigSGEcQgEpZdRQzAfpNUPHLYLTAnDHrtPTtYEYTxsPQzsIfUOm"
	)

	if ret := RandString(69); ret != first {
		t.Errorf("expected %s, got %s", first, ret)
	}

	if ret := RandString(69); ret != second {
		t.Errorf("expected %s, got %s", second, ret)
	}

	if ret := RandString(69); ret != third {
		t.Errorf("expected %s, got %s", third, ret)
	}

	t.Log("TestRandString passed")
}
