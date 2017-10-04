package notification

import "testing"

func TestNewSimpleHttpNotify(t *testing.T) {
	n := NewSimpleHttpNotify("test")
	if n == nil {
		t.Fatal("should not be nil")
	}
}
