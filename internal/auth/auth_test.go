package auth_test

import (
	"testing"
	"time"

	"github.com/Alody/Go-HTTP-server/internal/auth"
	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "supersecurepassword"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = auth.CheckPasswordHash(hash, password)
	if err != nil {
		t.Errorf("CheckPasswordHash failed: %v", err)
	}

	err = auth.CheckPasswordHash(hash, "wrongpassword")
	if err == nil {
		t.Error("CheckPasswordHash should have failed for wrong password, but it didn't")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "mysecretkey"
	userID := uuid.New()
	expiresIn := time.Minute

	token, err := auth.MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	returnedID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("ValidateJWT failed: %v", err)
	}
	if returnedID != userID {
		t.Errorf("Expected UUID %v, got %v", userID, returnedID)
	}
}

func TestValidateJWTWithInvalidSecret(t *testing.T) {
	secret := "correctsecret"
	wrongSecret := "wrongsecret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = auth.ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("ValidateJWT should have failed with wrong secret, but it didn't")
	}
}

func TestValidateExpiredJWT(t *testing.T) {
	secret := "secret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, -1*time.Minute) // already expired
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = auth.ValidateJWT(token, secret)
	if err == nil {
		t.Error("ValidateJWT should have failed for expired token, but it didn't")
	}
}
