package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestTokenRoundTrip(t *testing.T) {
	secret := []byte("test-secret")

	token, err := GenerateToken(secret, "user-123")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	userID, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if userID != "user-123" {
		t.Fatalf("expected user-123, got %q", userID)
	}
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	token, err := GenerateToken([]byte("secret-a"), "user-123")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	if _, err := ParseToken([]byte("secret-b"), token); err == nil {
		t.Fatal("expected error for token signed with a different secret")
	}
}

func TestParseTokenRejectsOtherAlgorithms(t *testing.T) {
	secret := []byte("test-secret")
	claims := tokenClaims{UserID: "user-123"}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(secret)
	if err != nil {
		t.Fatalf("SignedString: %v", err)
	}

	if _, err := ParseToken(secret, token); err == nil {
		t.Fatal("expected error for non-HS256 token")
	}
}

func TestDummyHashRejectsPasswords(t *testing.T) {
	if len(dummyHash) == 0 {
		t.Fatal("dummyHash was not generated")
	}
	if CheckPassword(string(dummyHash), "password123") {
		t.Fatal("dummyHash should not match an ordinary password")
	}
}
