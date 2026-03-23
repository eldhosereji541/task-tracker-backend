package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var testTokenSvc = NewTokenService([]byte("test-secret-key-that-is-at-least-32-chars!!"))

func TestGenerateToken_ReturnsNonEmpty(t *testing.T) {
	token, err := testTokenSvc.GenerateToken("user-555")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	token, err := testTokenSvc.GenerateToken("user-555")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := testTokenSvc.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if claims.UserID != "user-555" {
		t.Errorf("expected userID 'user-555', got '%s'", claims.UserID)
	}
}

func TestValidateToken_Tampered(t *testing.T) {
	token, _ := testTokenSvc.GenerateToken("user-555")

	_, err := testTokenSvc.ValidateToken(token + "tampered")
	if err == nil {
		t.Fatal("expected error for tampered token")
	}
}

func TestValidateToken_Empty(t *testing.T) {
	_, err := testTokenSvc.ValidateToken("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	svc := NewTokenService([]byte("test-secret-key-that-is-at-least-32-chars!!"))
	claims := &Claims{
		UserID: "user-555",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(svc.secret)

	_, err := svc.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	svc1 := NewTokenService([]byte("first-secret-key-that-is-32-chars-long!!"))
	svc2 := NewTokenService([]byte("different-secret-key-32-chars-long!!!!!!"))

	token, _ := svc1.GenerateToken("user-555")

	_, err := svc2.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error when validating with wrong secret")
	}
}
