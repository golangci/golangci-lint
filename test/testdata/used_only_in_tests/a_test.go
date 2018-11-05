package p

import "testing"

func TestF(t *testing.T) {
	if !f() {
		t.Fail()
	}
}
