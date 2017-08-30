package nativestore

import (
	"testing"
)

func TestSetGet(t *testing.T) {
	Set("localhost", "user", "password")
	user, secret, err := Get("localhost")
	if err != nil {
		t.Errorf("Should not get error, got %v", err)
	}

	if user != "user" {
		t.Errorf("Expecting user, got %s", user)
	}

	if secret != "password" {
		t.Errorf("Expecting password, got %s", secret)
	}
}
