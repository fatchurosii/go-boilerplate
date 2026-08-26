package auth

import (
	"testing"
	"time"
)

func TestJWTManagerGenerateAndParseAccessToken(t *testing.T) {
	manager, err := NewJWTManager("test-secret", "go-boilerplate", time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	token, err := manager.GenerateAccessToken("user-1")
	if err != nil {
		t.Fatal(err)
	}

	claims, err := manager.ParseAccessToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if claims.Subject != "user-1" {
		t.Fatalf("subject = %q, want %q", claims.Subject, "user-1")
	}
}
