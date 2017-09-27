package jwt

import (
	"log"
	"testing"
)

func TestNewCtx(t *testing.T) {
	ctx, err := NewCtx()
	if err != nil {
		t.Fatal(err)
	}

	if ctx == nil {
		t.Fatal("should not be nil")
	}
}

func TestGenerateToken(t *testing.T) {
	ctx, _ := NewCtx()
	token, stoken, _ := ctx.GenerateToken()
	if token == nil {
		t.Fatal("should not be nil")
	}

	log.Println(token, stoken)
}
