package auth
import (
	"golang.org/x/crypto/bcrypt"
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