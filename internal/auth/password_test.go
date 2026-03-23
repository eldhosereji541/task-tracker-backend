package auth

import "testing"

func TestHashPassword_ReturnsNonEmpty(t *testing.T) {
	hash, err := HashPassword("mypassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "mypassword" {
		t.Fatal("hash must not equal the plain password")
	}
}

func TestHashPassword_ProducesDifferentHashes(t *testing.T) {
	hash1, _ := HashPassword("samepassword")
	hash2, _ := HashPassword("samepassword")

	if hash1 == hash2 {
		t.Error("expected different hashes for same password (bcrypt must salt)")
	}
}

func TestCheckPasswordHash_Correct(t *testing.T) {
	hash, _ := HashPassword("correctpassword")

	if !CheckPasswordHash("correctpassword", hash) {
		t.Error("expected true for correct password")
	}
}

func TestCheckPasswordHash_Wrong(t *testing.T) {
	hash, _ := HashPassword("correctpassword")

	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("expected false for wrong password")
	}
}

func TestCheckPasswordHash_Empty(t *testing.T) {
	hash, _ := HashPassword("somepassword")

	if CheckPasswordHash("", hash) {
		t.Error("expected false for empty password")
	}
}
