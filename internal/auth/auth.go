package auth

import (
	"net/http"
	"time"

	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"crypto/rand"
	"encoding/hex"
)

// HashPassword takes a plain text password as input and returns a securely hashed version of it.
// This function uses bcrypt for hashing, which is suitable for storing passwords securely.
func HashPassword(password string) (string, error) {
	// Generate a hashed version of the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	// Compare the hashed password with the plain text password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err // Passwords do not match
	}
	return nil // Passwords match
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	},
	).SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return uuid.Parse(claims.Subject)
	}
	return uuid.Nil, jwt.NewValidationError("invalid token claims", jwt.ValidationErrorClaimsInvalid)
}

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if len(authHeader) == 0 {
		return "", jwt.NewValidationError("missing Authorization header", jwt.ValidationErrorMalformed)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", jwt.NewValidationError("missing Bearer token", jwt.ValidationErrorMalformed)
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if len(token) == 0 {
		return "", jwt.NewValidationError("empty token", jwt.ValidationErrorMalformed)
	}

	return token, nil
}

func MakeRefreshToken() (string, error) {
	
	token := make([]byte, 32) // 256-bit token
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil

}
