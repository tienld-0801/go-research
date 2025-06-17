package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Handle Compare Password
func CheckPassword(storedPassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	if err != nil {
		log.Println("Password comparison failed:", err)
		return false
	}
	return true
}

// Handle hash the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
