package testdata

import "testing"

// TestSimple is a basic test function
func TestSimple(t *testing.T) {
	if 1+1 != 2 {
		t.Error("math is broken")
	}
}
