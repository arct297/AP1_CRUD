package tools

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a plaintext password with a hashed password
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
