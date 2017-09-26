package jwt

import (
	"log"
	"testing"
)

func TestNewCtx(t *testing.T) {
	ctx, err := NewCtx("user")
	if err != nil {
		t.Fatal(err)
	}

	if ctx == nil {
		t.Fatal("should not be nil")
	}

	if ctx.User != "user" {
		t.Fatalf("expected user, got:", ctx.User)
	}

	log.Println("reuse:", ctx.Reuse)
}
